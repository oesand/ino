package ino

import (
	"context"
	"net/http"

	"github.com/oesand/ino/validate"
)

// RequestParam creates a ParamProvider that returns the entire HTTP request object.
// This is useful when you need access to the full request (headers, method, URL, etc.)
// in addition to or instead of individual parameters.
func RequestParam() ParamProvider[*http.Request] {
	return &requestParameter{}
}

type requestParameter struct{}

func (hp *requestParameter) Optional() ParamProvider[*http.Request] {
	return hp
}

func (hp *requestParameter) GetParamValue(request *http.Request) (*http.Request, validate.Errors) {
	return request, nil
}

// ContextParam creates a ParamProvider that returns the context from the HTTP request.
func ContextParam() ParamProvider[context.Context] {
	return &contextParameter{}
}

type contextParameter struct{}

func (hp *contextParameter) Optional() ParamProvider[context.Context] {
	return hp
}

func (hp *contextParameter) GetParamValue(request *http.Request) (context.Context, validate.Errors) {
	return request.Context(), nil
}
