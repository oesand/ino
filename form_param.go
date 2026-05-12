package ino

import (
	"fmt"
	"net/http"

	"github.com/oesand/ino/validate"
)

// FormParam creates a ParameterProvider that extracts a form parameter from the HTTP request.
// It supports validation and can be made optional.
func FormParam[T validate.BasicTypes](name string, validators ...validate.Validator[T]) ParameterProvider[T] {
	return &formParameter[T]{
		name:       name,
		validators: validators,
		post:       false,
	}
}

// PostFormParam creates a ParameterProvider that extracts a POST form parameter from the HTTP request.
// It supports validation and can be made optional.
func PostFormParam[T validate.BasicTypes](name string, validators ...validate.Validator[T]) ParameterProvider[T] {
	return &formParameter[T]{
		name:       name,
		validators: validators,
		post:       true,
	}
}

type formParameter[T validate.BasicTypes] struct {
	name       string
	optional   bool
	post       bool
	validators []validate.Validator[T]
}

func (fp *formParameter[T]) Optional() ParameterProvider[T] {
	fp.optional = true
	return fp
}

func (fp *formParameter[T]) GetParamValue(request *http.Request) (val T, errs validate.Errors) {
	var str string
	if fp.post {
		str = request.PostFormValue(fp.name)
	} else {
		str = request.FormValue(fp.name)
	}

	if str == "" {
		if !fp.optional {
			errs = []string{fmt.Sprintf("form param '%s' is required", fp.name)}
		}
		return
	}

	val, err := parseBasicTypes[T](str)
	if err != "" {
		errs = []string{fmt.Sprintf("form param '%s' %s", fp.name, err)}
		return
	}

	for _, validator := range fp.validators {
		for _, err := range validator.Validate(val) {
			errs = append(errs, fmt.Sprintf("form param '%s': %s", fp.name, err))
		}
	}
	return val, errs
}
