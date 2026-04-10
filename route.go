package ino

import (
	"fmt"
	"net/http"
)

// Route defines the interface for an HTTP route with method, pattern, handler, and attributes.
type Route interface {
	Method() string
	Pattern() string
	Handler() http.Handler
	Attrs() []any
}

// Get creates a GET route with the given pattern and handler.
func Get(pattern string, handler http.HandlerFunc, attrs ...any) Route {
	return Handle(http.MethodGet, pattern, handler, attrs...)
}

// Post creates a POST route with the given pattern and handler.
func Post(pattern string, handler http.HandlerFunc, attrs ...any) Route {
	return Handle(http.MethodPost, pattern, handler, attrs...)
}

// Put creates a PUT route with the given pattern and handler.
func Put(pattern string, handler http.HandlerFunc, attrs ...any) Route {
	return Handle(http.MethodPut, pattern, handler, attrs...)
}

// Delete creates a DELETE route with the given pattern and handler.
func Delete(pattern string, handler http.HandlerFunc, attrs ...any) Route {
	return Handle(http.MethodDelete, pattern, handler, attrs...)
}

// Options creates an OPTIONS route with the given pattern and handler.
func Options(pattern string, handler http.HandlerFunc, attrs ...any) Route {
	return Handle(http.MethodOptions, pattern, handler, attrs...)
}

// Head creates a HEAD route with the given pattern and handler.
func Head(pattern string, handler http.HandlerFunc, attrs ...any) Route {
	return Handle(http.MethodHead, pattern, handler, attrs...)
}

// Connect creates a CONNECT route with the given pattern and handler.
func Connect(pattern string, handler http.HandlerFunc, attrs ...any) Route {
	return Handle(http.MethodConnect, pattern, handler, attrs...)
}

// Patch creates a PATCH route with the given pattern and handler.
func Patch(pattern string, handler http.HandlerFunc, attrs ...any) Route {
	return Handle(http.MethodPatch, pattern, handler, attrs...)
}

// Trace creates a TRACE route with the given pattern and handler.
func Trace(pattern string, handler http.HandlerFunc, attrs ...any) Route {
	return Handle(http.MethodTrace, pattern, handler, attrs...)
}

// Handle creates a route with the given HTTP method, pattern, and handler.
// It validates that the pattern starts with '/', the method is valid, and the handler is not nil.
func Handle(method, pattern string, handler http.Handler, attrs ...any) Route {
	if pattern == "" {
		panic("mux: route pattern must have at least one character")
	}
	if pattern[0] != '/' {
		panic(fmt.Sprintf("mux: route pattern must starts with '/': %s", pattern))
	}

	if !IsValidMethod(method) {
		panic(fmt.Sprintf("mux: invalid http method: %s", pattern))
	}
	if handler == nil {
		panic(fmt.Sprintf("plow: nil handler: %s", pattern))
	}

	return &prefaceRoute{
		method:  method,
		pattern: pattern,
		handler: handler,
		attrs:   attrs,
	}
}

type prefaceRoute struct {
	method, pattern string
	handler         http.Handler
	attrs           []any
}

func (r *prefaceRoute) Method() string {
	return r.method
}

func (r *prefaceRoute) Pattern() string {
	return r.pattern
}

func (r *prefaceRoute) Handler() http.Handler {
	return r.handler
}

func (r *prefaceRoute) Attrs() []any {
	return r.attrs
}
