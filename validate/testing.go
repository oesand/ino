package validate

import (
	"testing"
)

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
