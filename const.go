package ino

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/oesand/ino/internal"
	"github.com/oesand/ino/validate"
)

var urlParamsKey = internal.CtxKey{Key: "mux/url_params"}
var matchedRouteKey = internal.CtxKey{Key: "mux/matched_route"}

// IsValidMethod checks whether the given HTTP method is valid.
func IsValidMethod(method string) bool {
	return method == http.MethodGet ||
		method == http.MethodPost ||
		method == http.MethodPut ||
		method == http.MethodDelete ||
		method == http.MethodOptions ||
		method == http.MethodHead ||
		method == http.MethodConnect ||
		method == http.MethodPatch ||
		method == http.MethodTrace
}

// UrlParams extracts URL parameters from the request context.
// Returns a map of parameter names to their values for the matched route.
func UrlParams(ctx context.Context) map[string]string {
	return ctx.Value(urlParamsKey).(map[string]string)
}

// MatchedRoute retrieves the matched route from the request context.
// Returns the Route that was matched during request routing.
func MatchedRoute(ctx context.Context) Route {
	return ctx.Value(matchedRouteKey).(Route)
}

func parseBasicTypes[T validate.BasicTypes](str string) (T, string) {
	var val T
	switch any(val).(type) {
	case string:
		val = any(str).(T)
	case bool:
		bv, err := strconv.ParseBool(str)
		if err != nil {
			return val, "must be bool"
		}
		val = any(bv).(T)
	case uint, uint8, uint16, uint32, uint64:
		bitSize := bitSizeNum(val)
		uiv, err := strconv.ParseUint(str, 10, bitSize)
		if err != nil {
			return val, "must be integer"
		}
		val = any(uiv).(T)
	case int, int8, int16, int32, int64:
		bitSize := bitSizeNum(val)
		iv, err := strconv.ParseInt(str, 10, bitSize)
		if err != nil {
			return val, "must be integer"
		}
		val = any(iv).(T)
	case float32, float64:
		bitSize := bitSizeNum(val)
		iv, err := strconv.ParseFloat(str, bitSize)
		if err != nil {
			return val, "must be float"
		}
		val = any(iv).(T)
	default:
		panic(fmt.Sprintf("ino: unknown type: %s", reflect.TypeFor[T]().String()))
	}
	return val, ""
}

func bitSizeNum[T any](v T) int {
	return int(unsafe.Sizeof(v) * 8)
}
