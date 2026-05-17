package validate_test

import (
	"testing"

	"github.com/oesand/ino/validate"
)

func TestFieldPrefixesError(t *testing.T) {
	type Parent struct{ Name string }

	desc := fieldDescriptor[Parent, string]{
		Name:  "Name",
		Value: func(p *Parent) string { return p.Name },
	}

	v := validate.Field(desc, validate.MinRunes(3))

	res := v.Validate(&Parent{Name: "ab"})
	if res.IsValid() {
		t.Fatalf("expected invalid")
	}
	if err := res.Error(); err != "'Name': must have at least 3 characters" {
		t.Fatalf("unexpected error %s", err)
	}
}

func TestFieldRPrefixesError(t *testing.T) {
	type Parent struct{ Name string }

	v := validate.FieldR[Parent, string]("Name", validate.MinRunes(3))

	res := v.Validate(&Parent{Name: "ab"})
	if res.IsValid() {
		t.Fatalf("expected invalid")
	}
	if err := res.Error(); err != "'Name': must have at least 3 characters" {
		t.Fatalf("unexpected error %s", err)
	}
}

func TestFieldRInvalidWithMessage(t *testing.T) {
	type Parent struct{ Name string }

	v := validate.FieldR[Parent, string]("Name", validate.MinRunes(3, validate.WithMessage("name too short")))

	res := v.Validate(&Parent{Name: "ab"})
	if res.IsValid() {
		t.Fatalf("expected invalid")
	}
	if err := res.Error(); err != "'Name': name too short" {
		t.Fatalf("unexpected error %s", err)
	}
}

func TestFieldRPanicsWhenFieldNotFound(t *testing.T) {
	type Parent struct{ Name string }

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for missing field")
		}
	}()

	_ = validate.FieldR[Parent, string]("Missing", validate.MinRunes(1))
}

type fieldDescriptor[Struct any, Field any] struct {
	Name  string
	Value func(*Struct) Field
}

func (desc fieldDescriptor[Struct, Field]) GetName() string {
	return desc.Name
}

func (desc fieldDescriptor[Struct, Field]) GetValue(s *Struct) Field {
	return desc.Value(s)
}
