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

// ParamHandler2 is a handler that takes 2 parameters.
func ParamHandler2[T1, T2 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	handler func(T1, T2, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, writer)
	}, param1, param2)
}

// ParamHandler3 is a handler that takes 3 parameters.
func ParamHandler3[T1, T2, T3 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	handler func(T1, T2, T3, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, writer)
	}, param1, param2, param3)
}

// ParamHandler4 is a handler that takes 4 parameters.
func ParamHandler4[T1, T2, T3, T4 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	handler func(T1, T2, T3, T4, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, writer)
	}, param1, param2, param3, param4)
}

// ParamHandler5 is a handler that takes 5 parameters.
func ParamHandler5[T1, T2, T3, T4, T5 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	handler func(T1, T2, T3, T4, T5, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, writer)
	}, param1, param2, param3, param4, param5)
}

// ParamHandler6 is a handler that takes 6 parameters.
func ParamHandler6[T1, T2, T3, T4, T5, T6 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	param6 ParamProvider[T6],
	handler func(T1, T2, T3, T4, T5, T6, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p6, paramErrs := param6.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, p6, writer)
	}, param1, param2, param3, param4, param5, param6)
}

// ParamHandler7 is a handler that takes 7 parameters.
func ParamHandler7[T1, T2, T3, T4, T5, T6, T7 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	param6 ParamProvider[T6],
	param7 ParamProvider[T7],
	handler func(T1, T2, T3, T4, T5, T6, T7, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p6, paramErrs := param6.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p7, paramErrs := param7.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, p6, p7, writer)
	}, param1, param2, param3, param4, param5, param6, param7)
}

// ParamHandler8 is a handler that takes 8 parameters.
func ParamHandler8[T1, T2, T3, T4, T5, T6, T7, T8 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	param6 ParamProvider[T6],
	param7 ParamProvider[T7],
	param8 ParamProvider[T8],
	handler func(T1, T2, T3, T4, T5, T6, T7, T8, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p6, paramErrs := param6.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p7, paramErrs := param7.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p8, paramErrs := param8.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, p6, p7, p8, writer)
	}, param1, param2, param3, param4, param5, param6, param7, param8)
}

// ParamHandler9 is a handler that takes 9 parameters.
func ParamHandler9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	param6 ParamProvider[T6],
	param7 ParamProvider[T7],
	param8 ParamProvider[T8],
	param9 ParamProvider[T9],
	handler func(T1, T2, T3, T4, T5, T6, T7, T8, T9, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p6, paramErrs := param6.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p7, paramErrs := param7.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p8, paramErrs := param8.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p9, paramErrs := param9.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, p6, p7, p8, p9, writer)
	}, param1, param2, param3, param4, param5, param6, param7, param8, param9)
}

// ParamHandler10 is a handler that takes 10 parameters.
func ParamHandler10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	param6 ParamProvider[T6],
	param7 ParamProvider[T7],
	param8 ParamProvider[T8],
	param9 ParamProvider[T9],
	param10 ParamProvider[T10],
	handler func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p6, paramErrs := param6.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p7, paramErrs := param7.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p8, paramErrs := param8.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p9, paramErrs := param9.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p10, paramErrs := param10.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, writer)
	}, param1, param2, param3, param4, param5, param6, param7, param8, param9, param10)
}

// ParamHandler11 is a handler that takes 11 parameters.
func ParamHandler11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	param6 ParamProvider[T6],
	param7 ParamProvider[T7],
	param8 ParamProvider[T8],
	param9 ParamProvider[T9],
	param10 ParamProvider[T10],
	param11 ParamProvider[T11],
	handler func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p6, paramErrs := param6.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p7, paramErrs := param7.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p8, paramErrs := param8.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p9, paramErrs := param9.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p10, paramErrs := param10.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p11, paramErrs := param11.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, writer)
	}, param1, param2, param3, param4, param5, param6, param7, param8, param9, param10, param11)
}

// ParamHandler12 is a handler that takes 12 parameters.
func ParamHandler12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	param6 ParamProvider[T6],
	param7 ParamProvider[T7],
	param8 ParamProvider[T8],
	param9 ParamProvider[T9],
	param10 ParamProvider[T10],
	param11 ParamProvider[T11],
	param12 ParamProvider[T12],
	handler func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p6, paramErrs := param6.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p7, paramErrs := param7.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p8, paramErrs := param8.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p9, paramErrs := param9.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p10, paramErrs := param10.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p11, paramErrs := param11.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p12, paramErrs := param12.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, writer)
	}, param1, param2, param3, param4, param5, param6, param7, param8, param9, param10, param11, param12)
}

// ParamHandler13 is a handler that takes 13 parameters.
func ParamHandler13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	param6 ParamProvider[T6],
	param7 ParamProvider[T7],
	param8 ParamProvider[T8],
	param9 ParamProvider[T9],
	param10 ParamProvider[T10],
	param11 ParamProvider[T11],
	param12 ParamProvider[T12],
	param13 ParamProvider[T13],
	handler func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p6, paramErrs := param6.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p7, paramErrs := param7.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p8, paramErrs := param8.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p9, paramErrs := param9.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p10, paramErrs := param10.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p11, paramErrs := param11.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p12, paramErrs := param12.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p13, paramErrs := param13.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, writer)
	}, param1, param2, param3, param4, param5, param6, param7, param8, param9, param10, param11, param12, param13)
}

// ParamHandler14 is a handler that takes 14 parameters.
func ParamHandler14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	param6 ParamProvider[T6],
	param7 ParamProvider[T7],
	param8 ParamProvider[T8],
	param9 ParamProvider[T9],
	param10 ParamProvider[T10],
	param11 ParamProvider[T11],
	param12 ParamProvider[T12],
	param13 ParamProvider[T13],
	param14 ParamProvider[T14],
	handler func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p6, paramErrs := param6.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p7, paramErrs := param7.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p8, paramErrs := param8.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p9, paramErrs := param9.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p10, paramErrs := param10.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p11, paramErrs := param11.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p12, paramErrs := param12.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p13, paramErrs := param13.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p14, paramErrs := param14.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, writer)
	}, param1, param2, param3, param4, param5, param6, param7, param8, param9, param10, param11, param12, param13, param14)
}

// ParamHandler15 is a handler that takes 15 parameters.
func ParamHandler15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	param6 ParamProvider[T6],
	param7 ParamProvider[T7],
	param8 ParamProvider[T8],
	param9 ParamProvider[T9],
	param10 ParamProvider[T10],
	param11 ParamProvider[T11],
	param12 ParamProvider[T12],
	param13 ParamProvider[T13],
	param14 ParamProvider[T14],
	param15 ParamProvider[T15],
	handler func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p6, paramErrs := param6.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p7, paramErrs := param7.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p8, paramErrs := param8.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p9, paramErrs := param9.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p10, paramErrs := param10.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p11, paramErrs := param11.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p12, paramErrs := param12.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p13, paramErrs := param13.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p14, paramErrs := param14.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p15, paramErrs := param15.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15, writer)
	}, param1, param2, param3, param4, param5, param6, param7, param8, param9, param10, param11, param12, param13, param14, param15)
}

// ParamHandler16 is a handler that takes 16 parameters.
func ParamHandler16[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	param6 ParamProvider[T6],
	param7 ParamProvider[T7],
	param8 ParamProvider[T8],
	param9 ParamProvider[T9],
	param10 ParamProvider[T10],
	param11 ParamProvider[T11],
	param12 ParamProvider[T12],
	param13 ParamProvider[T13],
	param14 ParamProvider[T14],
	param15 ParamProvider[T15],
	param16 ParamProvider[T16],
	handler func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p6, paramErrs := param6.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p7, paramErrs := param7.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p8, paramErrs := param8.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p9, paramErrs := param9.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p10, paramErrs := param10.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p11, paramErrs := param11.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p12, paramErrs := param12.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p13, paramErrs := param13.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p14, paramErrs := param14.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p15, paramErrs := param15.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p16, paramErrs := param16.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15, p16, writer)
	}, param1, param2, param3, param4, param5, param6, param7, param8, param9, param10, param11, param12, param13, param14, param15, param16)
}

// ParamHandler17 is a handler that takes 17 parameters.
func ParamHandler17[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	param6 ParamProvider[T6],
	param7 ParamProvider[T7],
	param8 ParamProvider[T8],
	param9 ParamProvider[T9],
	param10 ParamProvider[T10],
	param11 ParamProvider[T11],
	param12 ParamProvider[T12],
	param13 ParamProvider[T13],
	param14 ParamProvider[T14],
	param15 ParamProvider[T15],
	param16 ParamProvider[T16],
	param17 ParamProvider[T17],
	handler func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p6, paramErrs := param6.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p7, paramErrs := param7.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p8, paramErrs := param8.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p9, paramErrs := param9.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p10, paramErrs := param10.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p11, paramErrs := param11.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p12, paramErrs := param12.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p13, paramErrs := param13.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p14, paramErrs := param14.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p15, paramErrs := param15.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p16, paramErrs := param16.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p17, paramErrs := param17.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15, p16, p17, writer)
	}, param1, param2, param3, param4, param5, param6, param7, param8, param9, param10, param11, param12, param13, param14, param15, param16, param17)
}

// ParamHandler18 is a handler that takes 18 parameters.
func ParamHandler18[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	param6 ParamProvider[T6],
	param7 ParamProvider[T7],
	param8 ParamProvider[T8],
	param9 ParamProvider[T9],
	param10 ParamProvider[T10],
	param11 ParamProvider[T11],
	param12 ParamProvider[T12],
	param13 ParamProvider[T13],
	param14 ParamProvider[T14],
	param15 ParamProvider[T15],
	param16 ParamProvider[T16],
	param17 ParamProvider[T17],
	param18 ParamProvider[T18],
	handler func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p6, paramErrs := param6.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p7, paramErrs := param7.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p8, paramErrs := param8.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p9, paramErrs := param9.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p10, paramErrs := param10.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p11, paramErrs := param11.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p12, paramErrs := param12.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p13, paramErrs := param13.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p14, paramErrs := param14.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p15, paramErrs := param15.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p16, paramErrs := param16.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p17, paramErrs := param17.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p18, paramErrs := param18.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15, p16, p17, p18, writer)
	}, param1, param2, param3, param4, param5, param6, param7, param8, param9, param10, param11, param12, param13, param14, param15, param16, param17, param18)
}

// ParamHandler19 is a handler that takes 19 parameters.
func ParamHandler19[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	param6 ParamProvider[T6],
	param7 ParamProvider[T7],
	param8 ParamProvider[T8],
	param9 ParamProvider[T9],
	param10 ParamProvider[T10],
	param11 ParamProvider[T11],
	param12 ParamProvider[T12],
	param13 ParamProvider[T13],
	param14 ParamProvider[T14],
	param15 ParamProvider[T15],
	param16 ParamProvider[T16],
	param17 ParamProvider[T17],
	param18 ParamProvider[T18],
	param19 ParamProvider[T19],
	handler func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p6, paramErrs := param6.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p7, paramErrs := param7.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p8, paramErrs := param8.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p9, paramErrs := param9.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p10, paramErrs := param10.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p11, paramErrs := param11.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p12, paramErrs := param12.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p13, paramErrs := param13.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p14, paramErrs := param14.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p15, paramErrs := param15.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p16, paramErrs := param16.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p17, paramErrs := param17.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p18, paramErrs := param18.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p19, paramErrs := param19.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15, p16, p17, p18, p19, writer)
	}, param1, param2, param3, param4, param5, param6, param7, param8, param9, param10, param11, param12, param13, param14, param15, param16, param17, param18, param19)
}

// ParamHandler20 is a handler that takes 20 parameters.
func ParamHandler20[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, T20 any](
	param1 ParamProvider[T1],
	param2 ParamProvider[T2],
	param3 ParamProvider[T3],
	param4 ParamProvider[T4],
	param5 ParamProvider[T5],
	param6 ParamProvider[T6],
	param7 ParamProvider[T7],
	param8 ParamProvider[T8],
	param9 ParamProvider[T9],
	param10 ParamProvider[T10],
	param11 ParamProvider[T11],
	param12 ParamProvider[T12],
	param13 ParamProvider[T13],
	param14 ParamProvider[T14],
	param15 ParamProvider[T15],
	param16 ParamProvider[T16],
	param17 ParamProvider[T17],
	param18 ParamProvider[T18],
	param19 ParamProvider[T19],
	param20 ParamProvider[T20],
	handler func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, T20, http.ResponseWriter),
) http.Handler {
	return newParamHandler(func(writer http.ResponseWriter, request *http.Request) {
		var errs validate.Errors
		p1, paramErrs := param1.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p2, paramErrs := param2.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p3, paramErrs := param3.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p4, paramErrs := param4.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p5, paramErrs := param5.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p6, paramErrs := param6.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p7, paramErrs := param7.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p8, paramErrs := param8.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p9, paramErrs := param9.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p10, paramErrs := param10.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p11, paramErrs := param11.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p12, paramErrs := param12.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p13, paramErrs := param13.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p14, paramErrs := param14.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p15, paramErrs := param15.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p16, paramErrs := param16.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p17, paramErrs := param17.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p18, paramErrs := param18.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p19, paramErrs := param19.GetParamValue(request)
		errs = append(errs, paramErrs...)
		p20, paramErrs := param20.GetParamValue(request)
		errs = append(errs, paramErrs...)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15, p16, p17, p18, p19, p20, writer)
	}, param1, param2, param3, param4, param5, param6, param7, param8, param9, param10, param11, param12, param13, param14, param15, param16, param17, param18, param19, param20)
}
