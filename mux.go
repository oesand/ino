package ino

import (
	"container/list"
	"context"
	"fmt"
	"iter"
	"maps"
	"net/http"
	"sort"
)

// Middleware is a function type that wraps an HTTP handler with additional behavior.
type Middleware func(writer http.ResponseWriter, request *http.Request, next http.Handler)

// Mux is an HTTP request multiplexer that routes requests to registered handlers based on method and path pattern.
type Mux struct {
	routes      map[string][]*muxRoute
	middlewares *list.List
	notFound    http.Handler
}

type muxRoute struct {
	RoutePattern
	method  string
	handler http.Handler
	attrs   []any
}

func (m *muxRoute) Method() string {
	return m.method
}

func (m *muxRoute) Pattern() string {
	return m.Original
}

func (m *muxRoute) Handler() http.Handler {
	return m.handler
}

func (m *muxRoute) Attrs() []any {
	return m.attrs
}

// New creates a new Mux with optional initial routes.
func New(routes ...Route) *Mux {
	mux := &Mux{
		routes:      make(map[string][]*muxRoute),
		middlewares: list.New(),
	}
	mux.Include(routes...)
	return mux
}

// Include registers additional routes to the Mux.
func (mux *Mux) Include(routes ...Route) {
	if mux.routes == nil {
		mux.routes = make(map[string][]*muxRoute)
	}

	for _, route := range routes {
		pattern := route.Pattern()
		compiledPattern, err := ParseRoutePattern(pattern)
		if err != nil {
			panic(fmt.Sprintf("cannot compile route pattern: %s", pattern))
		}

		method := route.Method()
		mux.routes[method] = append(mux.routes[method], &muxRoute{
			RoutePattern: *compiledPattern,
			method:       method,
			handler:      route.Handler(),
			attrs:        route.Attrs(),
		})
	}

	for _, muxRoutes := range mux.routes {
		sort.Slice(muxRoutes, func(i, j int) bool {
			if muxRoutes[i].Depth == muxRoutes[j].Depth {
				return len(muxRoutes[i].ParamNames) < len(muxRoutes[j].ParamNames)
			}
			return muxRoutes[i].Depth > muxRoutes[j].Depth
		})
	}
}

// Routes returns an iterator over all registered routes.
func (mux *Mux) Routes() iter.Seq[Route] {
	return func(yield func(Route) bool) {
		for _, routes := range mux.routes {
			for _, route := range routes {
				if !yield(route) {
					return
				}
			}
		}
	}
}

// NotFound sets the handler to be called when no route matches the request.
func (mux *Mux) NotFound(handler http.Handler) {
	mux.notFound = handler
}

// Middleware registers a middleware to be applied to all requests.
func (mux *Mux) Middleware(middleware Middleware) {
	mux.middlewares.PushBack(middleware)
}

// ServeHTTP implements the http.Handler interface, routing the request to the appropriate handler.
func (mux *Mux) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var route Route
	var urlParams map[string]string

	routes, found := mux.routes[request.Method]
	if found {
		found = false
		for _, r := range routes {
			matched, params := r.Match(request.URL.Path)
			if !matched {
				continue
			}

			route = r
			urlParams = maps.Collect(params)
			found = true
			break
		}
	}

	if found {
		ctx := request.Context()

		ctx = context.WithValue(ctx, matchedRouteKey, route)

		if len(urlParams) > 0 {
			ctx = context.WithValue(ctx, urlParamsKey, urlParams)
		}

		request = request.WithContext(ctx)
	}

	handler := muxHandler{
		nextMiddleware: mux.middlewares.Front(),
		notFound:       mux.notFound,
		route:          route,
	}
	handler.ServeHTTP(writer, request)
}

type muxHandler struct {
	nextMiddleware *list.Element
	notFound       http.Handler
	route          Route
}

func (m *muxHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if m.nextMiddleware != nil {
		middleware := m.nextMiddleware.Value.(Middleware)
		m.nextMiddleware = m.nextMiddleware.Next()
		middleware(writer, request, m)
		return
	}

	if m.route != nil {
		m.route.Handler().ServeHTTP(writer, request)
		return
	}

	if nf := m.notFound; nf != nil {
		nf.ServeHTTP(writer, request)
		return
	}

	http.NotFound(writer, request)
}
