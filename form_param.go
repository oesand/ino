package ino

import (
	"fmt"
	"net/http"

	"github.com/oesand/ino/validate"
)

func FormParam[T validate.BasicTypes](name string, validators ...validate.Validator[T]) ParameterProvider[T] {
	return &formParameter[T]{
		name:       name,
		validators: validators,
	}
}

type formParameter[T validate.BasicTypes] struct {
	name       string
	optional   bool
	validators []validate.Validator[T]
}

func (fp *formParameter[T]) Optional() ParameterProvider[T] {
	fp.optional = true
	return fp
}

func (fp *formParameter[T]) GetParamValue(request *http.Request) (val T, errs validate.Errors) {
	str := request.FormValue(fp.name)

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
