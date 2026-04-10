package ino

import (
	"fmt"
	"net/http"

	"github.com/oesand/ino/validate"
)

// UrlParam creates a ParameterProvider that extracts a URL parameter from the request.
// The parameter value is retrieved from the URL parameters stored in the request context
// (set by the router during route matching). It supports basic types like string, int, bool, etc.
// Optional validators can be provided to validate the parameter value.
// Use Optional() on the returned provider to make the parameter optional.
func UrlParam[T validate.BasicTypes](name string, validators ...validate.Validator[T]) ParameterProvider[T] {
	return &urlParameter[T]{
		name:       name,
		validators: validators,
	}
}

type urlParameter[T validate.BasicTypes] struct {
	name       string
	optional   bool
	validators []validate.Validator[T]
}

func (up *urlParameter[T]) Optional() ParameterProvider[T] {
	up.optional = true
	return up
}

func (up *urlParameter[T]) GetParamValue(request *http.Request) (val T, errs validate.Errors) {
	urlParams := UrlParams(request.Context())
	if urlParams == nil {
		if !up.optional {
			errs = []string{fmt.Sprintf("url param '%s' is required", up.name)}
		}
		return
	}

	str, _ := urlParams[up.name]
	if str == "" {
		if !up.optional {
			errs = []string{fmt.Sprintf("url param '%s' is required", up.name)}
		}
		return
	}

	val, err := parseBasicTypes[T](str)
	if err != "" {
		errs = []string{fmt.Sprintf("url param '%s' %s", up.name, err)}
		return
	}

	for _, validator := range up.validators {
		for _, err := range validator.Validate(val) {
			errs = append(errs, fmt.Sprintf("url param '%s': %s", up.name, err))
		}
	}
	return val, errs
}
