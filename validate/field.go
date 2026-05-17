package validate

import (
	"fmt"
	"reflect"
)

// Field returns a validator for a struct field described by `descriptor`.
// The returned validator extracts the field value from the parent struct and
// runs the provided validators; any errors are prefixed with the field name.
func Field[Struct any, Field any](descriptor FieldDescriptor[Struct, Field], validators ...Validator[Field]) Validator[*Struct] {
	return &fieldValidator[Struct, Field]{
		fieldName:   descriptor.GetName(),
		fieldGetter: descriptor.GetValue,
		validators:  validators,
	}
}

// FieldR returns a validator for a struct field selected by name.
//
// The returned validator locates the named field using reflection, extracts its
// value, and runs the provided validators. Any validation errors are prefixed
// with the field name.
func FieldR[Struct any, Field any](name string, validators ...Validator[Field]) Validator[*Struct] {
	field, has := reflect.TypeFor[Struct]().FieldByName(name)
	if !has {
		panic(fmt.Sprintf("field %s not found", name))
	}

	return &fieldValidator[Struct, Field]{
		fieldName: name,
		fieldGetter: func(s *Struct) Field {
			val := reflect.ValueOf(s).Elem().FieldByIndex(field.Index).Interface()
			if val == nil {
				var zero Field
				return zero
			}
			return val.(Field)
		},
		validators: validators,
	}
}

type fieldValidator[Struct any, Field any] struct {
	fieldName   string
	fieldGetter func(*Struct) Field
	validators  []Validator[Field]
}

func (validator *fieldValidator[Struct, Field]) Validate(parent *Struct) Errors {
	var errors []string
	name := validator.fieldName
	value := validator.fieldGetter(parent)
	for _, v := range validator.validators {
		for _, err := range v.Validate(value) {
			errors = append(errors, fmt.Sprintf("'%s': %s", name, err))
		}
	}
	return errors
}
