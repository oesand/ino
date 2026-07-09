package openapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/oesand/ino"
	"github.com/oesand/ino/collection"
	"github.com/oesand/ino/internal"
)

// newOperation creates the OpenAPI path pattern and operation wrapper for a route.
func newOperation(route ino.Route) (string, *operationWrapper) {
	method := route.Method()
	pattern := route.Pattern()
	if compiledRoute, ok := route.(internal.CompiledRoute); ok {
		pattern = compiledRoute.ClearedPattern()
	}
	operation := operationWrapper(spec.Operation{
		OperationProps: spec.OperationProps{
			ID:      generateOperationID(method, pattern),
			Summary: fmt.Sprintf("%s %s", method, pattern),
			Responses: &spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					StatusCodeResponses: make(map[int]spec.Response),
				},
			},
		},
	})

	return pattern, &operation
}

type operationWrapper spec.Operation

// Base returns the wrapped OpenAPI operation.
func (wrapper *operationWrapper) Base() *spec.Operation {
	return (*spec.Operation)(wrapper)
}

// Apply assigns the wrapped operation to the matching HTTP method on a path item.
func (wrapper *operationWrapper) Apply(pathItem *spec.PathItem, method string) {
	operation := wrapper.Base()
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

// FillConsumes adds parameters and consumed content types inferred from the route handler.
func (wrapper *operationWrapper) FillConsumes(route ino.Route, swagger *swaggerWrapper) {
	var declaredPathParams collection.Set[string]
	if paramHandler, ok := route.Handler().(internal.ParamHandler); ok {
		var consumes collection.Set[string]

		isUpdate := route.Method() == http.MethodPost || route.Method() == http.MethodPatch ||
			route.Method() == http.MethodPut
		for _, schema := range paramHandler.Params() {
			param, consume := swagger.MakeParamFromHandlerParam(schema, isUpdate)
			if param != nil {
				wrapper.Parameters = append(wrapper.Parameters, *param)
			}
			if consume != "" {
				consumes.Add(consume)
			}
			if schema.ParamType() == internal.PathParamType {
				declaredPathParams.Add(schema.Name())
			}
		}

		wrapper.Consumes = consumes.Values()
	}

	if compiledRoute, ok := route.(internal.CompiledRoute); ok {
		for _, name := range compiledRoute.PathParams() {
			if declaredPathParams.Has(name) {
				continue
			}

			param := spec.PathParam(name).Typed("string", "")
			wrapper.Parameters = append(wrapper.Parameters, *param)
			declaredPathParams.Add(name)
		}
	}
}

// ScanAttributes applies route attributes to responses, tags, security, parameters, and produced content types.
func (wrapper *operationWrapper) ScanAttributes(route ino.Route, swagger *swaggerWrapper) {
	var produces collection.Set[string]
	for _, attr := range route.Attrs() {
		switch at := attr.(type) {
		case *ResponseAttr:
			if at.ContentType != "" {
				produces.Add(at.ContentType)
			}

			schema := at.Schema
			if schema == nil && at.Type != nil {
				schema = swagger.MakeSchema(at.Type)
			}

			wrapper.Responses.StatusCodeResponses[at.Code] = spec.Response{
				ResponseProps: spec.ResponseProps{
					Schema:      schema,
					Description: at.Description,
				},
			}
		case *TagAttr:
			wrapper.Tags = append(wrapper.Tags, at.Name)
		case *SummaryAttr:
			wrapper.Summary = at.Summary
		case *AuthAttr:
			wrapper.Security = append(wrapper.Security, map[string][]string{
				at.Name: at.Scopes,
			})
		case *spec.Parameter:
			wrapper.Parameters = append(wrapper.Parameters, *at)
		}
	}
	wrapper.Produces = produces.Values()
}

// generateOperationID derives a stable operation ID from the HTTP method and path pattern.
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
