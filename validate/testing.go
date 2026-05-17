package validate

import (
	"fmt"
	"reflect"
	"testing"
)

// Must validate the given value using the provided validators and calls t.Error
// if any validation errors occur. It collects all errors from all validators.
func Must[Value any](t *testing.T, value Value, validators ...Validator[Value]) {
	var errors Errors
	for _, v := range validators {
		errors = append(errors, v.Validate(value)...)
	}
	if len(errors) > 0 {
		t.Error(errors.Error())
	}
}

// DeepEqual returns a validator that checks if the actual value is deeply equal
// to the expected value using reflect.DeepEqual.
func DeepEqual[Value any](expected Value, options ...func(setter optionsSetters)) Validator[Value] {
	validator := &deepEqualizer[Value]{expected: expected}

	for _, option := range options {
		option(validator)
	}

	return validator
}

type deepEqualizer[Value any] struct {
	expected     Value
	errorMessage string
}

func (de *deepEqualizer[Value]) SetErrorMessage(err string) {
	de.errorMessage = err
}

func (de *deepEqualizer[Value]) Validate(actual Value) Errors {
	if !reflect.DeepEqual(actual, de.expected) {
		errorMessage := de.errorMessage
		if errorMessage == "" {
			errorMessage = fmt.Sprintf("value expected %+v, but got %+v", de.expected, actual)
		}
		return Errors{errorMessage}
	}
	return nil
}
