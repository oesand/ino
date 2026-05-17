package validate

import (
	"fmt"
)

// Slice returns a validator that applies the provided element validators to
// each element of a slice. If any element produces validation errors the
// returned result will contain those errors prefixed with the element index.
func Slice[Element any](validators ...Validator[Element]) Validator[[]Element] {
	return &sliceValidator[Element]{
		validators: validators,
	}
}

type sliceValidator[Element any] struct {
	validators []Validator[Element]
}

func (validator *sliceValidator[Element]) Validate(slice []Element) Errors {
	var errors []string
	for i, el := range slice {
		for _, validator := range validator.validators {
			for _, err := range validator.Validate(el) {
				errors = append(errors, fmt.Sprintf("[%d]: %s", i, err))
			}
		}
		if len(errors) > 0 {
			break
		}
	}

	return errors
}

// MinCount returns a validator that ensures a slice has at least `minLength`
// elements.
func MinCount[Element any](minLength int, options ...func(setter optionsSetters)) Validator[[]Element] {
	validator := &sliceMinValidator[Element]{minLength: minLength}

	for _, option := range options {
		option(validator)
	}

	return validator
}

type sliceMinValidator[Element any] struct {
	minLength    int
	errorMessage string
}

func (validator *sliceMinValidator[Value]) SetErrorMessage(err string) {
	validator.errorMessage = err
}

func (validator *sliceMinValidator[Element]) Validate(slice []Element) Errors {
	if len(slice) < validator.minLength {
		errorMessage := validator.errorMessage
		if errorMessage == "" {
			errorMessage = fmt.Sprintf("count must be greater than or equal to %v", validator.minLength)
		}
		return []string{errorMessage}
	}
	return nil
}

// MaxCount returns a validator that ensures a slice has at most `maxLength`
// elements.
func MaxCount[Element any](maxLength int, options ...func(setter optionsSetters)) Validator[[]Element] {
	validator := &sliceMaxValidator[Element]{maxLength: maxLength}

	for _, option := range options {
		option(validator)
	}

	return validator
}

type sliceMaxValidator[Element any] struct {
	maxLength    int
	errorMessage string
}

func (validator *sliceMaxValidator[Value]) SetErrorMessage(err string) {
	validator.errorMessage = err
}

func (validator *sliceMaxValidator[Element]) Validate(slice []Element) Errors {
	if len(slice) > validator.maxLength {
		errorMessage := validator.errorMessage
		if errorMessage == "" {
			errorMessage = fmt.Sprintf("count must be less than or equal to %v", validator.maxLength)
		}
		return []string{errorMessage}
	}
	return nil
}
