package validate

import (
	"fmt"
)

// Field returns a validator for a struct field described by `descriptor`.
// The returned validator extracts the field value from the parent struct and
// runs the provided validators; any errors are prefixed with the field name.
func Field[Struct any, Field any](descriptor FieldDescriptor[Struct, Field], validators ...Validator[Field]) Validator[*Struct] {
	return &fieldValidator[Struct, Field]{
		descriptor: descriptor,
		validators: validators,
	}
}

type fieldValidator[Struct any, Field any] struct {
	descriptor FieldDescriptor[Struct, Field]
	validators []Validator[Field]
}

func (validator *fieldValidator[Struct, Field]) Validate(parent *Struct) Errors {
	var errors []string
	name := validator.descriptor.GetName()
	value := validator.descriptor.GetValue(parent)
	for _, v := range validator.validators {
		for _, err := range v.Validate(value) {
			errors = append(errors, fmt.Sprintf("> '%s': %s", name, err))
		}
	}
	return errors
}
