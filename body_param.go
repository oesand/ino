package mo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/oesand/mo/internal"
	"github.com/oesand/mo/validate"
)

// BodyParam creates a ParamProvider that returns the raw request body as an io.ReadCloser.
// This gives direct access to the request body stream, useful for streaming large files,
// custom parsing, or when you don't want the framework to buffer the entire body.
func BodyParam() ParamProvider[io.ReadCloser] {
	return &bodyParameter{}
}

var _ internal.ParamSchema = (*bodyParameter)(nil)

type bodyParameter struct {
	optional bool
}

func (param *bodyParameter) Name() string {
	return ""
}

func (param *bodyParameter) ParamType() internal.ParamType {
	return internal.RawBodyParamType
}

func (param *bodyParameter) Type() reflect.Type {
	return nil
}

func (param *bodyParameter) IsRequired() bool {
	return !param.optional
}

func (param *bodyParameter) Optional() ParamProvider[io.ReadCloser] {
	param.optional = true
	return param
}

func (param *bodyParameter) GetParamValue(request *http.Request) (io.ReadCloser, validate.Errors) {
	if request.Body == nil {
		if !param.optional {
			return nil, []string{"body is required"}
		}
	}
	return request.Body, nil
}

// JsonParam creates a ParamProvider that parses JSON from the request body into a struct.
// The JSON is decoded into the specified type T, and optional validators can be applied
// to the parsed object.
func JsonParam[T any](validators ...validate.Validator[*T]) ParamProvider[*T] {
	return &jsonParameter[T]{validators: validators}
}

var _ internal.ParamSchema = (*jsonParameter[struct{}])(nil)

type jsonParameter[T any] struct {
	optional   bool
	validators []validate.Validator[*T]
}

func (param *jsonParameter[T]) Name() string {
	return ""
}

func (param *jsonParameter[T]) ParamType() internal.ParamType {
	return internal.JsonBodyParamType
}

func (param *jsonParameter[T]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}

func (param *jsonParameter[T]) IsRequired() bool {
	return !param.optional
}

func (param *jsonParameter[T]) Optional() ParamProvider[*T] {
	param.optional = true
	return param
}

func (param *jsonParameter[T]) GetParamValue(request *http.Request) (*T, validate.Errors) {
	if request.Body == nil {
		if !param.optional {
			return nil, []string{"json body is required"}
		}
		return nil, nil
	}

	var value T
	err := json.NewDecoder(request.Body).Decode(&value)
	if err != nil {
		if !param.optional {
			return nil, []string{"json body is required"}
		}
		return nil, nil
	}

	var errs []string
	for _, validator := range param.validators {
		for _, err := range validator.Validate(&value) {
			errs = append(errs, fmt.Sprintf("json body: %s", err))
		}
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return &value, errs
}
