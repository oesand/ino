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
	if err := res.Error(); err != "> 'Name': must have at least 3 characters" {
		t.Fatalf("unexpected error %s", err)
	}
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
