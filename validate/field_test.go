package validate_test

import (
	"testing"

	"github.com/oesand/ino/validate"
	"github.com/oesand/octo/octogen"
)

func TestFieldPrefixesError(t *testing.T) {
	type Parent struct{ Name string }

	desc := octogen.FieldDescriptor[Parent, string]{
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
