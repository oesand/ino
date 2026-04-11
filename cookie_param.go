package ino

import (
	"fmt"
	"net/http"

	"github.com/oesand/ino/validate"
)

// CookieParam creates a ParameterProvider that extracts and validates a cookie value from an HTTP request.
// The cookie value is parsed into the specified type T (string, int64, or bool) and validated
// using the provided validators. If no validators are provided, only type parsing is performed.
func CookieParam[T validate.BasicTypes](name string, validators ...validate.Validator[T]) ParameterProvider[T] {
	return &cookieParameter[T]{
		name:       name,
		validators: validators,
	}
}

type cookieParameter[T validate.BasicTypes] struct {
	name       string
	optional   bool
	validators []validate.Validator[T]
}

func (hp *cookieParameter[T]) Optional() ParameterProvider[T] {
	hp.optional = true
	return hp
}

func (hp *cookieParameter[T]) GetParamValue(request *http.Request) (val T, errs validate.Errors) {
	cookie, err := request.Cookie(hp.name)
	if err != nil {
		if !hp.optional {
			errs = []string{fmt.Sprintf("cookie '%s' is required", hp.name)}
		}
		return
	}
	str := cookie.Value

	val, perr := parseBasicTypes[T](str)
	if perr != "" {
		errs = []string{fmt.Sprintf("cookie '%s' %s", hp.name, perr)}
		return
	}

	for _, validator := range hp.validators {
		for _, err := range validator.Validate(val) {
			errs = append(errs, fmt.Sprintf("cookie '%s': %s", hp.name, err))
		}
	}
	return val, errs
}
