package mo

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/oesand/mo/internal"
	"github.com/oesand/mo/validate"
)

// PathParam creates a ParamProvider that extracts a URL parameter from the request.
// The parameter value is retrieved from the URL parameters stored in the request context
// (set by the router during route matching). It supports basic types like string, int, bool, etc.
// Optional validators can be provided to validate the parameter value.
// Use Optional() on the returned provider to make the parameter optional.
func PathParam[T validate.BasicTypes](name string, validators ...validate.Validator[T]) ParamProvider[T] {
	return &pathParameter[T]{
		name:       name,
		validators: validators,
	}
}

var _ internal.ParamSchema = (*pathParameter[string])(nil)

type pathParameter[T validate.BasicTypes] struct {
	name       string
	optional   bool
	validators []validate.Validator[T]
}

func (param *pathParameter[T]) Name() string {
	return param.name
}

func (param *pathParameter[T]) ParamType() internal.ParamType {
	return internal.PathParamType
}

func (param *pathParameter[T]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}

func (param *pathParameter[T]) IsRequired() bool {
	return !param.optional
}

func (param *pathParameter[T]) Optional() ParamProvider[T] {
	param.optional = true
	return param
}

func (param *pathParameter[T]) GetParamValue(request *http.Request) (val T, errs validate.Errors) {
	pathParams := PathParams(request.Context())
	if pathParams == nil {
		if !param.optional {
			errs = []string{fmt.Sprintf("path param '%s' is required", param.name)}
		}
		return
	}

	str, _ := pathParams[param.name]
	if str == "" {
		if !param.optional {
			errs = []string{fmt.Sprintf("path param '%s' is required", param.name)}
		}
		return
	}

	val, err := parseBasicTypes[T](str)
	if err != "" {
		errs = []string{fmt.Sprintf("path param '%s' %s", param.name, err)}
		return
	}

	for _, validator := range param.validators {
		for _, err := range validator.Validate(val) {
			errs = append(errs, fmt.Sprintf("path param '%s': %s", param.name, err))
		}
	}
	return val, errs
}
