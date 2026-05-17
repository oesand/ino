package validate

import (
	"fmt"

	"github.com/oesand/ino/collection"
)

// OneOf returns a `Validator` that checks whether a value
// is one of the provided `values`.
func OneOf[Element comparable](values []Element, options ...func(setter optionsSetters)) Validator[Element] {
	var valuesString string
	for i, value := range values {
		if i > 0 {
			valuesString += fmt.Sprintf(", %v", value)
		} else {
			valuesString += fmt.Sprint(value)
		}
	}

	validator := &oneOfValidator[Element]{
		values:       collection.SetOf(values...),
		valuesString: valuesString,
	}

	for _, option := range options {
		option(validator)
	}

	return validator
}

type oneOfValidator[Element comparable] struct {
	values       collection.Set[Element]
	valuesString string
	errorMessage string
}

func (validator *oneOfValidator[Value]) SetErrorMessage(err string) {
	validator.errorMessage = err
}

func (validator *oneOfValidator[Element]) Validate(value Element) Errors {
	if !validator.values.Has(value) {
		errorMessage := validator.errorMessage
		if errorMessage == "" {
			errorMessage = fmt.Sprintf("must be in %s", validator.valuesString)
		}
		return []string{errorMessage}
	}
	return nil
}
