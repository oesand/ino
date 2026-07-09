package validate_test

import (
	"testing"

	"github.com/oesand/mo/validate"
)

func TestRule_Valid(t *testing.T) {
	v := validate.Rule(func(value int) bool {
		return value > 0
	})

	err := v.Validate(5)
	if !err.IsValid() {
		t.Fatalf("expected valid, got errors: %v", err)
	}
}

func TestRule_InvalidDefaultMessage(t *testing.T) {
	v := validate.Rule(func(value int) bool {
		return value < 0
	})

	err := v.Validate(1)
	if err.IsValid() {
		t.Fatal("expected invalid, got valid")
	}
	if err.Error() != "defined rule was not met" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestRule_InvalidWithMessage(t *testing.T) {
	v := validate.Rule(func(value int) bool {
		return value%2 == 0
	}, validate.WithMessage("value must be even"))

	err := v.Validate(3)
	if err.IsValid() {
		t.Fatal("expected invalid, got valid")
	}
	if err.Error() != "value must be even" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}
