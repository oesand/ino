package ino

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/oesand/ino/internal"
	"github.com/oesand/ino/validate"
)

// HeaderParam creates a ParamProvider that extracts an HTTP header from the request.
// The header value is retrieved using request.Header.Get(). It supports basic types like string, int, bool, etc.
// Optional validators can be provided to validate the header value.
// Use Optional() on the returned provider to make the header optional.
func HeaderParam[T validate.BasicTypes](name string, validators ...validate.Validator[T]) ParamProvider[T] {
	return &headerParameter[T]{
		name:       name,
		validators: validators,
	}
}

var _ internal.ParamSchema = (*headerParameter[string])(nil)

type headerParameter[T validate.BasicTypes] struct {
	name       string
	optional   bool
	validators []validate.Validator[T]
}

func (param *headerParameter[T]) Name() string {
	return param.name
}

func (param *headerParameter[T]) ParamType() internal.ParamType {
	return internal.HeaderParamType
}

func (param *headerParameter[T]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}

func (param *headerParameter[T]) IsRequired() bool {
	return !param.optional
}

func (param *headerParameter[T]) Optional() ParamProvider[T] {
	param.optional = true
	return param
}

func (param *headerParameter[T]) GetParamValue(request *http.Request) (value T, errs validate.Errors) {
	str := request.Header.Get(param.name)
	if str == "" {
		if !param.optional {
			errs = []string{fmt.Sprintf("header '%s' is required", param.name)}
		}
		return
	}

	val, err := parseBasicTypes[T](str)
	if err != "" {
		errs = []string{fmt.Sprintf("header '%s' %s", param.name, err)}
		return
	}

	for _, validator := range param.validators {
		for _, err := range validator.Validate(val) {
			errs = append(errs, fmt.Sprintf("header '%s': %s", param.name, err))
		}
	}
	return val, errs
}
