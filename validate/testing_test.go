package validate_test

import (
	"testing"

	"github.com/oesand/ino/validate"
)

func TestMust_NoErrors(t *testing.T) {
	// Test that Must does not call t.Error when there are no validation errors
	validate.Must(t, 5, validate.DeepEqual(5))
	// If the test passes, it means no t.Error was called
}

func TestDeepEqual_Equal(t *testing.T) {
	v := validate.DeepEqual(5)
	errs := v.Validate(5)
	if !errs.IsValid() {
		t.Errorf("expected valid, got errors: %v", errs)
	}
}

func TestDeepEqual_NotEqual(t *testing.T) {
	v := validate.DeepEqual(5)
	errs := v.Validate(10)
	if errs.IsValid() {
		t.Error("expected invalid, got valid")
	}
	if errs.Error() != "value expected 5 but got 10" {
		t.Errorf("unexpected error: %s", errs.Error())
	}
}

func TestDeepEqual_Slice(t *testing.T) {
	expected := []int{1, 2, 3}
	v := validate.DeepEqual(expected)
	errs := v.Validate([]int{1, 2, 3})
	if !errs.IsValid() {
		t.Errorf("expected valid, got errors: %v", errs)
	}
	errs = v.Validate([]int{1, 2, 4})
	if errs.IsValid() {
		t.Error("expected invalid, got valid")
	}
}
