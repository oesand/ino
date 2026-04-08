package mux

import (
	"context"
	"net/http"

	"github.com/oesand/ino/internal"
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
