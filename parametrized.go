package mo

import (
	"net/http"

	"github.com/oesand/mo/internal"
	"github.com/oesand/mo/validate"
)

// ParamProvider is a provider that can be used to get a parameter value.
type ParamProvider[T any] interface {
	GetParamValue(*http.Request) (T, validate.Errors)
	Optional() ParamProvider[T]
}

func newParamHandler(handler http.HandlerFunc, parameters ...any) *paramHandler {
	var schemas []internal.ParamSchema
	for _, param := range parameters {
		if schema, ok := param.(internal.ParamSchema); ok {
			schemas = append(schemas, schema)
		}
	}
	return &paramHandler{
		handler:    handler,
		parameters: schemas,
	}
}

var _ internal.ParamHandler = (*paramHandler)(nil)

type paramHandler struct {
	handler    http.HandlerFunc
	parameters []internal.ParamSchema
}

func (handler *paramHandler) Params() []internal.ParamSchema {
	return handler.parameters
}

func (handler *paramHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler.handler(writer, request)
}

// ParamHandler1 is a handler that takes a single parameter.
func ParamHandler1[T1 any](
	param1 ParamProvider[T1],
	handler func(T1, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		p1, errs := param1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, writer)
	}, param1)
}
