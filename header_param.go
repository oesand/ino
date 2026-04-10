package ino

import (
	"fmt"
	"net/http"

	"github.com/oesand/ino/validate"
)

// HeaderParam creates a ParameterProvider that extracts an HTTP header from the request.
// The header value is retrieved using request.Header.Get(). It supports basic types like string, int, bool, etc.
// Optional validators can be provided to validate the header value.
// Use Optional() on the returned provider to make the header optional.
func HeaderParam[T validate.BasicTypes](name string, validators ...validate.Validator[T]) ParameterProvider[T] {
	return &headerParameter[T]{
		name:       name,
		validators: validators,
	}
}

type headerParameter[T validate.BasicTypes] struct {
	name       string
	optional   bool
	validators []validate.Validator[T]
}

func (hp *headerParameter[T]) Optional() ParameterProvider[T] {
	hp.optional = true
	return hp
}

func (hp *headerParameter[T]) GetParamValue(request *http.Request) (value T, errs validate.Errors) {
	str := request.Header.Get(hp.name)
	if str == "" {
		if !hp.optional {
			errs = []string{fmt.Sprintf("header '%s' is required", hp.name)}
		}
		return
	}

	val, err := parseBasicTypes[T](str)
	if err != "" {
		errs = []string{fmt.Sprintf("header '%s' %s", hp.name, err)}
		return
	}

	for _, validator := range hp.validators {
		for _, err := range validator.Validate(val) {
			errs = append(errs, fmt.Sprintf("header '%s': %s", hp.name, err))
		}
	}
	return val, errs
}
