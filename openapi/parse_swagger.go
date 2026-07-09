package openapi

import (
	"reflect"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/oesand/ino/internal"
)

// newSwagger creates a Swagger 2.0 wrapper with initialized paths and definitions.
func newSwagger(info *spec.Info) *swaggerWrapper {
	swagger := swaggerWrapper(spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger:     "2.0",
			Info:        info,
			Paths:       &spec.Paths{Paths: map[string]spec.PathItem{}},
			Definitions: spec.Definitions{},
		},
	})

	return &swagger
}

type swaggerWrapper spec.Swagger

// Base returns the wrapped Swagger document.
func (wrapper *swaggerWrapper) Base() *spec.Swagger {
	return (*spec.Swagger)(wrapper)
}

// MakeParamFromHandlerParam converts an ino parameter schema into an OpenAPI parameter.
func (wrapper *swaggerWrapper) MakeParamFromHandlerParam(paramSchema internal.ParamSchema, isUpdate bool) (*spec.Parameter, string) {
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
		if !isUpdate {
			parameter = spec.QueryParam(paramSchema.Name())
			break
		}
		fallthrough
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
		schema := wrapper.MakeSchema(paramType)
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

// MakeSchema returns a schema for typ and registers named struct definitions by reference.
func (wrapper *swaggerWrapper) MakeSchema(typ reflect.Type) *spec.Schema {
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

		if _, found := wrapper.Definitions[refName]; !found {
			wrapper.Definitions[refName] = *structToSchema(typ)
		}
		return spec.RefSchema("#/definitions/" + refName)
	}
	return TypeToSchema(typ)
}
