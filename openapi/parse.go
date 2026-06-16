package openapi

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/oesand/ino"
	"github.com/oesand/ino/collection"
	"github.com/oesand/ino/internal"
)

type GenerateOptions struct {
	Info spec.InfoProps
	Tags []spec.Tag
}

func GenerateSchema(mux *ino.Mux, options *GenerateOptions) (*spec.Swagger, error) {
	if mux == nil {
		return nil, errors.New("ino: mux is required")
	}

	var info spec.Info

	if options != nil {
		info.InfoProps = options.Info
	}

	if info.Title == "" {
		info.Title = "API"
	}

	if info.Version == "" {
		info.Version = "1.0.0"
	}

	swagger := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger:     "2.0",
			Info:        &info,
			Paths:       &spec.Paths{Paths: map[string]spec.PathItem{}},
			Definitions: spec.Definitions{},
		},
	}

	if options != nil {
		swagger.Tags = options.Tags
	}

	for route := range mux.Routes() {
		path := route.Pattern()
		pathItem := swagger.Paths.Paths[path]
		operation := &spec.Operation{
			OperationProps: spec.OperationProps{
				ID:      generateOperationID(route.Method(), path),
				Summary: fmt.Sprintf("%s %s", route.Method(), path),
				Responses: &spec.Responses{
					ResponsesProps: spec.ResponsesProps{
						StatusCodeResponses: make(map[int]spec.Response),
					},
				},
			},
		}

		var declaredPathParams collection.Set[string]
		if paramHandler, ok := route.Handler().(internal.ParamHandler); ok {
			var consumes collection.Set[string]

			isGetter := route.Method() == http.MethodGet || route.Method() == http.MethodHead
			for _, schema := range paramHandler.Params() {
				param, consume := parameterFromCustomSchema(schema, swagger.Definitions, isGetter)
				if param != nil {
					operation.Parameters = append(operation.Parameters, *param)
				}
				if consume != "" {
					consumes.Add(consume)
				}
				if schema.ParamType() == internal.PathParamType {
					declaredPathParams.Add(schema.Name())
				}
			}

			operation.Consumes = consumes.Values()
		}

		if compiledRoute, ok := route.(internal.CompiledRoute); ok {
			for _, name := range compiledRoute.PathParams() {
				if declaredPathParams.Has(name) {
					continue
				}

				param := spec.PathParam(name).Typed("string", "")
				operation.Parameters = append(operation.Parameters, *param)
				declaredPathParams.Add(name)
			}
		}

		var produces collection.Set[string]
		for _, attr := range route.Attrs() {
			switch at := attr.(type) {
			case *ResponseAttr:
				if at.ContentType != "" {
					produces.Add(at.ContentType)
				}

				schema := at.Schema
				if schema == nil && at.Type != nil {
					schema = generateAndDeclareSchema(at.Type, swagger.Definitions)
				}

				operation.Responses.StatusCodeResponses[at.Code] = spec.Response{
					ResponseProps: spec.ResponseProps{
						Schema:      schema,
						Description: at.Description,
					},
				}
			case *TagAttr:
				operation.Tags = append(operation.Tags, at.Name)
			}
		}
		operation.Produces = produces.Values()

		applyMethodOperation(&pathItem, route.Method(), operation)
		swagger.Paths.Paths[path] = pathItem
	}

	return swagger, nil
}

func applyMethodOperation(pathItem *spec.PathItem, method string, operation *spec.Operation) {
	switch method {
	case http.MethodGet:
		pathItem.Get = operation
	case http.MethodPost:
		pathItem.Post = operation
	case http.MethodPut:
		pathItem.Put = operation
	case http.MethodDelete:
		pathItem.Delete = operation
	case http.MethodPatch:
		pathItem.Patch = operation
	case http.MethodHead:
		pathItem.Head = operation
	case http.MethodOptions:
		pathItem.Options = operation
	default:
		pathItem.Get = operation
	}
}

func parameterFromCustomSchema(paramSchema internal.ParamSchema, definitions spec.Definitions, isGetter bool) (*spec.Parameter, string) {
	var parameter *spec.Parameter
	var consume string
	switch paramSchema.ParamType() {
	case internal.PathParamType:
		parameter = spec.PathParam(paramSchema.Name())
	case internal.HeaderParamType:
		parameter = spec.HeaderParam(paramSchema.Name())
	case internal.CookieParamType:
		parameter = &spec.Parameter{ParamProps: spec.ParamProps{Name: paramSchema.Name(), In: "cookie"}}
	case internal.FormParamType:
		if isGetter {
			parameter = spec.QueryParam(paramSchema.Name())
		} else {
			parameter = spec.FormDataParam(paramSchema.Name())
			consume = "application/x-www-form-urlencoded"
		}
	case internal.PostFormParamType:
		parameter = spec.FormDataParam(paramSchema.Name())
		consume = "application/x-www-form-urlencoded"
	case internal.FileParamType:
		parameter = spec.FileParam(paramSchema.Name())
		consume = "multipart/form-data"
		return parameter, consume
	case internal.MultipartFormParamType:
		return nil, "multipart/form-data"
	case internal.JsonBodyParamType:
		parameter = &spec.Parameter{ParamProps: spec.ParamProps{Name: "json", In: "body"}}
		consume = "application/json"
	case internal.RawBodyParamType:
		parameter = &spec.Parameter{
			ParamProps: spec.ParamProps{
				Name: "json", In: "body",
				Schema: new(spec.Schema).Typed("string", "byte"),
			},
		}
		return parameter, "application/octet-stream"
	default:
		parameter = spec.QueryParam(paramSchema.Name())
	}

	parameter.Required = paramSchema.IsRequired()

	if paramType := paramSchema.Type(); paramType != nil {
		schema := generateAndDeclareSchema(paramType, definitions)
		parameter.Schema = schema

		if types := schema.Type; len(types) > 0 {
			parameter.Type = types[0]
		}

		parameter.Format = schema.Format

		if parameter.Type == "array" && schema.Items != nil && schema.Items.Schema != nil {
			item := schema.Items.Schema
			parameter.Items = &spec.Items{SimpleSchema: spec.SimpleSchema{Type: item.Type[0], Format: item.Format}}
			parameter.CollectionFormat = "multi"
		}
	}

	return parameter, consume
}

func generateAndDeclareSchema(typ reflect.Type, definitions spec.Definitions) *spec.Schema {
	if typ == nil {
		return nil
	}
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	if typ.Kind() == reflect.Struct && typ.Name() != "" {
		refName := typ.Name()
		if pkg := typ.PkgPath(); pkg != "" {
			refName = strings.ReplaceAll(pkg, "/", "\\") + "." + refName
		}

		if _, found := definitions[refName]; !found {
			definitions[refName] = *structToSchema(reflect.New(typ).Elem().Interface())
		}
		return spec.RefSchema("#/definitions/" + refName)
	}
	return TypeToSchema(typ)
}

func generateOperationID(method, path string) string {
	result := strings.Trim(path, "/")
	result = strings.ReplaceAll(result, "{", "")
	result = strings.ReplaceAll(result, "}", "")
	result = strings.ReplaceAll(result, "/", "_")
	result = strings.ReplaceAll(result, "-", "_")
	if result == "" {
		result = "root"
	}
	return strings.ToLower(method) + "_" + result
}
