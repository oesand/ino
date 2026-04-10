package ino

import (
	"fmt"
)

// StackRoutes is a collection of routes that can be nested and prefixed.
type StackRoutes []Route

// Routes groups the given routes and attributes into a StackRoutes.
func Routes(options ...any) StackRoutes {
	return makeRoutes("", options)
}

// PrefixRoutes groups the given routes with a common path prefix and attributes into a StackRoutes.
func PrefixRoutes(prefix string, options ...any) StackRoutes {
	if len(prefix) < 2 {
		panic("mux: router prefix must have at least two characters")
	}
	if prefix[0] != '/' {
		panic(fmt.Sprintf("mux: router prefix must starts with '/': %s", prefix))
	}

	return makeRoutes(prefix, options)
}

func makeRoutes(prefix string, options []any) StackRoutes {
	var routes []Route
	var attrs []any
	for _, option := range options {
		switch val := option.(type) {
		case Route:
			routes = append(routes, val)
		case StackRoutes:
			routes = append(routes, val...)
		default:
			attrs = append(attrs, val)
		}
	}

	for i, route := range routes {
		pattern := route.Pattern()
		if prefix != "" {
			pattern = prefix + pattern
		}

		var finalAttrs []any
		finalAttrs = append(finalAttrs, route.Attrs()...)
		finalAttrs = append(finalAttrs, attrs...)
		routes[i] = Handle(route.Method(), pattern, route.Handler(), finalAttrs...)
	}

	return routes
}
