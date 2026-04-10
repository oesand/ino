package ino

import (
	"encoding/json"
	"net/http"

	"github.com/oesand/ino/validate"
)

// ParameterProvider is a provider that can be used to get a parameter value.
type ParameterProvider[T any] interface {
	GetParamValue(*http.Request) (T, validate.Errors)
	Optional() ParameterProvider[T]
}

func Errors(writer http.ResponseWriter, errors []string, code int) {
	header := writer.Header()

	header.Set("Content-Type", "application/json; charset=utf-8")
	header.Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(code)
	json.NewEncoder(writer).Encode(&errorResponse{Errors: errors})
}

type errorResponse struct {
	Errors []string `json:"errors"`
}

// ParamHandler is a handler that takes a single parameter.
func ParamHandler[T0 any](
	provider ParameterProvider[T0],
	handler func(T0, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, writer)
	}
}

// ParamHandler1 is a handler that takes two parameters.
func ParamHandler1[T0 any, T1 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	handler func(T0, T1, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, writer)
	}
}

// ParamHandler2 is a handler that takes three parameters.
func ParamHandler2[T0 any, T1 any, T2 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	handler func(T0, T1, T2, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, writer)
	}
}

// ParamHandler3 is a handler that takes four parameters.
func ParamHandler3[T0 any, T1 any, T2 any, T3 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	handler func(T0, T1, T2, T3, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, writer)
	}
}

// ParamHandler4 is a handler that takes five parameters.
func ParamHandler4[T0 any, T1 any, T2 any, T3 any, T4 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	handler func(T0, T1, T2, T3, T4, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, writer)
	}
}

// ParamHandler5 is a handler that takes six parameters.
func ParamHandler5[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	handler func(T0, T1, T2, T3, T4, T5, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, writer)
	}
}

// ParamHandler6 is a handler that takes seven parameters.
func ParamHandler6[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	provider6 ParameterProvider[T6],
	handler func(T0, T1, T2, T3, T4, T5, T6, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p6, errs := provider6.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, p6, writer)
	}
}

// ParamHandler7 is a handler that takes eight parameters.
func ParamHandler7[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	provider6 ParameterProvider[T6],
	provider7 ParameterProvider[T7],
	handler func(T0, T1, T2, T3, T4, T5, T6, T7, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p6, errs := provider6.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p7, errs := provider7.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, p6, p7, writer)
	}
}

// ParamHandler8 is a handler that takes nine parameters.
func ParamHandler8[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	provider6 ParameterProvider[T6],
	provider7 ParameterProvider[T7],
	provider8 ParameterProvider[T8],
	handler func(T0, T1, T2, T3, T4, T5, T6, T7, T8, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p6, errs := provider6.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p7, errs := provider7.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p8, errs := provider8.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, p6, p7, p8, writer)
	}
}

// ParamHandler9 is a handler that takes ten parameters.
func ParamHandler9[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	provider6 ParameterProvider[T6],
	provider7 ParameterProvider[T7],
	provider8 ParameterProvider[T8],
	provider9 ParameterProvider[T9],
	handler func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p6, errs := provider6.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p7, errs := provider7.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p8, errs := provider8.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p9, errs := provider9.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, writer)
	}
}

// ParamHandler10 is a handler that takes eleven parameters.
func ParamHandler10[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any, T10 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	provider6 ParameterProvider[T6],
	provider7 ParameterProvider[T7],
	provider8 ParameterProvider[T8],
	provider9 ParameterProvider[T9],
	provider10 ParameterProvider[T10],
	handler func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p6, errs := provider6.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p7, errs := provider7.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p8, errs := provider8.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p9, errs := provider9.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p10, errs := provider10.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, writer)
	}
}

// ParamHandler11 is a handler that takes twelve parameters.
func ParamHandler11[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any, T10 any, T11 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	provider6 ParameterProvider[T6],
	provider7 ParameterProvider[T7],
	provider8 ParameterProvider[T8],
	provider9 ParameterProvider[T9],
	provider10 ParameterProvider[T10],
	provider11 ParameterProvider[T11],
	handler func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p6, errs := provider6.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p7, errs := provider7.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p8, errs := provider8.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p9, errs := provider9.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p10, errs := provider10.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p11, errs := provider11.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, writer)
	}
}

// ParamHandler12 is a handler that takes thirteen parameters.
func ParamHandler12[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any, T10 any, T11 any, T12 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	provider6 ParameterProvider[T6],
	provider7 ParameterProvider[T7],
	provider8 ParameterProvider[T8],
	provider9 ParameterProvider[T9],
	provider10 ParameterProvider[T10],
	provider11 ParameterProvider[T11],
	provider12 ParameterProvider[T12],
	handler func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p6, errs := provider6.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p7, errs := provider7.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p8, errs := provider8.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p9, errs := provider9.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p10, errs := provider10.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p11, errs := provider11.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p12, errs := provider12.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, writer)
	}
}

// ParamHandler13 is a handler that takes fourteen parameters.
func ParamHandler13[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any, T10 any, T11 any, T12 any, T13 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	provider6 ParameterProvider[T6],
	provider7 ParameterProvider[T7],
	provider8 ParameterProvider[T8],
	provider9 ParameterProvider[T9],
	provider10 ParameterProvider[T10],
	provider11 ParameterProvider[T11],
	provider12 ParameterProvider[T12],
	provider13 ParameterProvider[T13],
	handler func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p6, errs := provider6.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p7, errs := provider7.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p8, errs := provider8.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p9, errs := provider9.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p10, errs := provider10.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p11, errs := provider11.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p12, errs := provider12.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p13, errs := provider13.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, writer)
	}
}

// ParamHandler14 is a handler that takes fifteen parameters.
func ParamHandler14[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any, T10 any, T11 any, T12 any, T13 any, T14 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	provider6 ParameterProvider[T6],
	provider7 ParameterProvider[T7],
	provider8 ParameterProvider[T8],
	provider9 ParameterProvider[T9],
	provider10 ParameterProvider[T10],
	provider11 ParameterProvider[T11],
	provider12 ParameterProvider[T12],
	provider13 ParameterProvider[T13],
	provider14 ParameterProvider[T14],
	handler func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p6, errs := provider6.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p7, errs := provider7.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p8, errs := provider8.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p9, errs := provider9.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p10, errs := provider10.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p11, errs := provider11.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p12, errs := provider12.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p13, errs := provider13.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p14, errs := provider14.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, writer)
	}
}

// ParamHandler15 is a handler that takes sixteen parameters.
func ParamHandler15[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any, T10 any, T11 any, T12 any, T13 any, T14 any, T15 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	provider6 ParameterProvider[T6],
	provider7 ParameterProvider[T7],
	provider8 ParameterProvider[T8],
	provider9 ParameterProvider[T9],
	provider10 ParameterProvider[T10],
	provider11 ParameterProvider[T11],
	provider12 ParameterProvider[T12],
	provider13 ParameterProvider[T13],
	provider14 ParameterProvider[T14],
	provider15 ParameterProvider[T15],
	handler func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p6, errs := provider6.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p7, errs := provider7.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p8, errs := provider8.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p9, errs := provider9.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p10, errs := provider10.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p11, errs := provider11.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p12, errs := provider12.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p13, errs := provider13.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p14, errs := provider14.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p15, errs := provider15.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15, writer)
	}
}

// ParamHandler16 is a handler that takes seventeen parameters.
func ParamHandler16[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any, T10 any, T11 any, T12 any, T13 any, T14 any, T15 any, T16 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	provider6 ParameterProvider[T6],
	provider7 ParameterProvider[T7],
	provider8 ParameterProvider[T8],
	provider9 ParameterProvider[T9],
	provider10 ParameterProvider[T10],
	provider11 ParameterProvider[T11],
	provider12 ParameterProvider[T12],
	provider13 ParameterProvider[T13],
	provider14 ParameterProvider[T14],
	provider15 ParameterProvider[T15],
	provider16 ParameterProvider[T16],
	handler func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p6, errs := provider6.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p7, errs := provider7.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p8, errs := provider8.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p9, errs := provider9.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p10, errs := provider10.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p11, errs := provider11.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p12, errs := provider12.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p13, errs := provider13.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p14, errs := provider14.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p15, errs := provider15.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p16, errs := provider16.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15, p16, writer)
	}
}

// ParamHandler17 is a handler that takes eighteen parameters.
func ParamHandler17[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any, T10 any, T11 any, T12 any, T13 any, T14 any, T15 any, T16 any, T17 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	provider6 ParameterProvider[T6],
	provider7 ParameterProvider[T7],
	provider8 ParameterProvider[T8],
	provider9 ParameterProvider[T9],
	provider10 ParameterProvider[T10],
	provider11 ParameterProvider[T11],
	provider12 ParameterProvider[T12],
	provider13 ParameterProvider[T13],
	provider14 ParameterProvider[T14],
	provider15 ParameterProvider[T15],
	provider16 ParameterProvider[T16],
	provider17 ParameterProvider[T17],
	handler func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p6, errs := provider6.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p7, errs := provider7.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p8, errs := provider8.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p9, errs := provider9.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p10, errs := provider10.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p11, errs := provider11.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p12, errs := provider12.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p13, errs := provider13.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p14, errs := provider14.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p15, errs := provider15.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p16, errs := provider16.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p17, errs := provider17.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15, p16, p17, writer)
	}
}

// ParamHandler18 is a handler that takes nineteen parameters.
func ParamHandler18[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any, T10 any, T11 any, T12 any, T13 any, T14 any, T15 any, T16 any, T17 any, T18 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	provider6 ParameterProvider[T6],
	provider7 ParameterProvider[T7],
	provider8 ParameterProvider[T8],
	provider9 ParameterProvider[T9],
	provider10 ParameterProvider[T10],
	provider11 ParameterProvider[T11],
	provider12 ParameterProvider[T12],
	provider13 ParameterProvider[T13],
	provider14 ParameterProvider[T14],
	provider15 ParameterProvider[T15],
	provider16 ParameterProvider[T16],
	provider17 ParameterProvider[T17],
	provider18 ParameterProvider[T18],
	handler func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p6, errs := provider6.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p7, errs := provider7.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p8, errs := provider8.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p9, errs := provider9.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p10, errs := provider10.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p11, errs := provider11.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p12, errs := provider12.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p13, errs := provider13.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p14, errs := provider14.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p15, errs := provider15.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p16, errs := provider16.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p17, errs := provider17.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p18, errs := provider18.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15, p16, p17, p18, writer)
	}
}

// ParamHandler19 is a handler that takes twenty parameters.
func ParamHandler19[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any, T10 any, T11 any, T12 any, T13 any, T14 any, T15 any, T16 any, T17 any, T18 any, T19 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	provider6 ParameterProvider[T6],
	provider7 ParameterProvider[T7],
	provider8 ParameterProvider[T8],
	provider9 ParameterProvider[T9],
	provider10 ParameterProvider[T10],
	provider11 ParameterProvider[T11],
	provider12 ParameterProvider[T12],
	provider13 ParameterProvider[T13],
	provider14 ParameterProvider[T14],
	provider15 ParameterProvider[T15],
	provider16 ParameterProvider[T16],
	provider17 ParameterProvider[T17],
	provider18 ParameterProvider[T18],
	provider19 ParameterProvider[T19],
	handler func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p6, errs := provider6.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p7, errs := provider7.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p8, errs := provider8.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p9, errs := provider9.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p10, errs := provider10.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p11, errs := provider11.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p12, errs := provider12.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p13, errs := provider13.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p14, errs := provider14.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p15, errs := provider15.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p16, errs := provider16.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p17, errs := provider17.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p18, errs := provider18.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p19, errs := provider19.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15, p16, p17, p18, p19, writer)
	}
}

// ParamHandler20 is a handler that takes twenty-one parameters.
func ParamHandler20[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any, T10 any, T11 any, T12 any, T13 any, T14 any, T15 any, T16 any, T17 any, T18 any, T19 any, T20 any](
	provider ParameterProvider[T0],
	provider1 ParameterProvider[T1],
	provider2 ParameterProvider[T2],
	provider3 ParameterProvider[T3],
	provider4 ParameterProvider[T4],
	provider5 ParameterProvider[T5],
	provider6 ParameterProvider[T6],
	provider7 ParameterProvider[T7],
	provider8 ParameterProvider[T8],
	provider9 ParameterProvider[T9],
	provider10 ParameterProvider[T10],
	provider11 ParameterProvider[T11],
	provider12 ParameterProvider[T12],
	provider13 ParameterProvider[T13],
	provider14 ParameterProvider[T14],
	provider15 ParameterProvider[T15],
	provider16 ParameterProvider[T16],
	provider17 ParameterProvider[T17],
	provider18 ParameterProvider[T18],
	provider19 ParameterProvider[T19],
	provider20 ParameterProvider[T20],
	handler func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T18, T19, T20, http.ResponseWriter),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		p0, errs := provider.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p1, errs := provider1.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p2, errs := provider2.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p3, errs := provider3.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p4, errs := provider4.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p5, errs := provider5.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p6, errs := provider6.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p7, errs := provider7.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p8, errs := provider8.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p9, errs := provider9.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p10, errs := provider10.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p11, errs := provider11.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p12, errs := provider12.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p13, errs := provider13.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p14, errs := provider14.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p15, errs := provider15.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p16, errs := provider16.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p17, errs := provider17.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p18, errs := provider18.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p19, errs := provider19.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		p20, errs := provider20.GetParamValue(request)
		if len(errs) > 0 {
			Errors(writer, errs, http.StatusUnprocessableEntity)
			return
		}
		handler(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15, p16, p17, p18, p19, p20, writer)
	}
}
