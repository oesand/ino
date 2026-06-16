package ino

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/oesand/ino/internal"
	"github.com/oesand/ino/validate"
)

// CookieParam creates a ParamProvider that extracts and validates a cookie value from an HTTP request.
// The cookie value is parsed into the specified type T (string, int64, or bool) and validated
// using the provided validators. If no validators are provided, only type parsing is performed.
func CookieParam[T validate.BasicTypes](name string, validators ...validate.Validator[T]) ParamProvider[T] {
	return &cookieParameter[T]{
		name:       name,
		validators: validators,
	}
}

var _ internal.ParamSchema = (*cookieParameter[string])(nil)

type cookieParameter[T validate.BasicTypes] struct {
	name       string
	optional   bool
	validators []validate.Validator[T]
}

func (param *cookieParameter[T]) Name() string {
	return param.name
}

func (param *cookieParameter[T]) ParamType() internal.ParamType {
	return internal.CookieParamType
}

func (param *cookieParameter[T]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}

func (param *cookieParameter[T]) IsRequired() bool {
	return !param.optional
}

func (param *cookieParameter[T]) Optional() ParamProvider[T] {
	param.optional = true
	return param
}

func (param *cookieParameter[T]) GetParamValue(request *http.Request) (val T, errs validate.Errors) {
	cookie, err := request.Cookie(param.name)
	if err != nil {
		if !param.optional {
			errs = []string{fmt.Sprintf("cookie '%s' is required", param.name)}
		}
		return
	}
	str := cookie.Value

	val, perr := parseBasicTypes[T](str)
	if perr != "" {
		errs = []string{fmt.Sprintf("cookie '%s' %s", param.name, perr)}
		return
	}

	for _, validator := range param.validators {
		for _, err := range validator.Validate(val) {
			errs = append(errs, fmt.Sprintf("cookie '%s': %s", param.name, err))
		}
	}
	return val, errs
}
