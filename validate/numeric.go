package validate

import (
	"fmt"
)

// Min creates a condition that validates a numeric value is greater than or equal to the minimum.
func Min[Value NumericTypes](min Value, options ...func(setter optionsSetters)) Validator[Value] {
	validator := &numericMinValidator[Value]{min: min}

	for _, option := range options {
		option(validator)
	}

	return validator
}

type numericMinValidator[Value NumericTypes] struct {
	min          Value
	errorMessage string
}

func (validator *numericMinValidator[Value]) SetErrorMessage(err string) {
	validator.errorMessage = err
}

func (validator *numericMinValidator[Value]) Validate(value Value) Errors {
	if value < validator.min {
		errorMessage := validator.errorMessage
		if errorMessage == "" {
			errorMessage = fmt.Sprintf("must be greater than or equal to %v", validator.min)
		}
		return []string{errorMessage}
	}
	return nil
}

// Max creates a condition that validates a numeric value is less than or equal to the maximum.
func Max[Value NumericTypes](max Value, options ...func(setter optionsSetters)) Validator[Value] {
	validator := &numericMaxValidator[Value]{max: max}

	for _, option := range options {
		option(validator)
	}

	return validator
}

type numericMaxValidator[Value NumericTypes] struct {
	max          Value
	errorMessage string
}

func (validator *numericMaxValidator[Value]) SetErrorMessage(err string) {
	validator.errorMessage = err
}

func (validator *numericMaxValidator[Value]) Validate(value Value) Errors {
	if value > validator.max {
		errorMessage := validator.errorMessage
		if errorMessage == "" {
			errorMessage = fmt.Sprintf("must be less than or equal to %v", validator.max)
		}
		return []string{errorMessage}
	}
	return nil
}
