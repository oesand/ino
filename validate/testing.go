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
		for _, err := range v.Validate(value) {
			errors = append(errors, err)
		}
	}
	if len(errors) != 0 {
		t.Error(errors.Error())
	}
}

// DeepEqual returns a validator that checks if the actual value is deeply equal
// to the expected value using reflect.DeepEqual.
func DeepEqual[Value any](expected Value) Validator[Value] {
	return &deepEqualizer[Value]{
		expected: expected,
	}
}

type deepEqualizer[Value any] struct {
	expected Value
}

func (de *deepEqualizer[Value]) Validate(actual Value) Errors {
	if !reflect.DeepEqual(actual, de.expected) {
		return Errors{fmt.Sprintf("value expected %v but got %v", de.expected, actual)}
	}
	return nil
}
