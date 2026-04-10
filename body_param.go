package ino

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/oesand/ino/validate"
)

// RequestParam creates a ParameterProvider that returns the entire HTTP request object.
// This is useful when you need access to the full request (headers, method, URL, etc.)
// in addition to or instead of individual parameters.
func RequestParam() ParameterProvider[*http.Request] {
	return &requestParameter{}
}

type requestParameter struct{}

func (hp *requestParameter) Optional() ParameterProvider[*http.Request] {
	return hp
}

func (hp *requestParameter) GetParamValue(request *http.Request) (*http.Request, validate.Errors) {
	return request, nil
}

// BodyParam creates a ParameterProvider that returns the raw request body as an io.ReadCloser.
// This gives direct access to the request body stream, useful for streaming large files,
// custom parsing, or when you don't want the framework to buffer the entire body.
func BodyParam() ParameterProvider[io.ReadCloser] {
	return &bodyParameter{}
}

type bodyParameter struct {
	optional bool
}

func (bp *bodyParameter) Optional() ParameterProvider[io.ReadCloser] {
	bp.optional = true
	return bp
}

func (bp *bodyParameter) GetParamValue(request *http.Request) (io.ReadCloser, validate.Errors) {
	if request.Body == nil {
		if !bp.optional {
			return nil, []string{"body is required"}
		}
	}
	return request.Body, nil
}

// MultipartFormParam creates a ParameterProvider that parses and returns multipart form data.
// The maxMemory parameter controls how much of the form data is stored in memory before
// spilling to temporary files on disk.
func MultipartFormParam(maxMemory int64) ParameterProvider[*multipart.Form] {
	return &multipartFormParameter{
		maxMemory: maxMemory,
	}
}

type multipartFormParameter struct {
	maxMemory int64
	optional  bool
}

func (mpp *multipartFormParameter) Optional() ParameterProvider[*multipart.Form] {
	mpp.optional = true
	return mpp
}

func (mpp *multipartFormParameter) GetParamValue(request *http.Request) (*multipart.Form, validate.Errors) {
	err := request.ParseMultipartForm(mpp.maxMemory)
	if err != nil {
		if !mpp.optional {
			return nil, []string{"multipart form is required"}
		}
		return nil, nil
	}
	return request.MultipartForm, nil
}

// JsonParam creates a ParameterProvider that parses JSON from the request body into a struct.
// The JSON is decoded into the specified type T, and optional validators can be applied
// to the parsed object.
func JsonParam[T any](validators ...validate.Validator[*T]) ParameterProvider[*T] {
	return &jsonParameter[T]{
		validators: validators,
	}
}

type jsonParameter[T any] struct {
	optional   bool
	validators []validate.Validator[*T]
}

func (jp *jsonParameter[T]) Optional() ParameterProvider[*T] {
	jp.optional = true
	return jp
}

func (jp *jsonParameter[T]) GetParamValue(request *http.Request) (*T, validate.Errors) {
	if request.Body == nil {
		if !jp.optional {
			return nil, []string{"json body is required"}
		}
		return nil, nil
	}

	var value T
	err := json.NewDecoder(request.Body).Decode(&value)
	if err != nil {
		if !jp.optional {
			return nil, []string{"json body is required"}
		}
		return nil, nil
	}

	var errs []string
	for _, validator := range jp.validators {
		for _, err := range validator.Validate(&value) {
			errs = append(errs, fmt.Sprintf("json body: %s", err))
		}
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return &value, errs
}
