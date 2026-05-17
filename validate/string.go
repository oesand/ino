package validate

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

// Regex returns a validator that checks the given string matches the provided
// regular expression.
//
// It panics if `regex` is nil.
func Regex(regex *regexp.Regexp, options ...func(setter optionsSetters)) Validator[string] {
	if regex == nil {
		panic("validate: regex is nil")
	}
	validator := &stringRegexValidator{regex: regex}

	for _, option := range options {
		option(validator)
	}

	return validator
}

type stringRegexValidator struct {
	regex        *regexp.Regexp
	errorMessage string
}

func (validator *stringRegexValidator) SetErrorMessage(err string) {
	validator.errorMessage = err
}

func (validator *stringRegexValidator) Validate(value string) Errors {
	if !validator.regex.MatchString(value) {
		errorMessage := validator.errorMessage
		if errorMessage == "" {
			errorMessage = "mismatch expected pattern"
		}
		return []string{errorMessage}
	}
	return nil
}

// RunesExactly returns a validator that ensures the string contains exactly
// `length` runes (Unicode code points).
func RunesExactly(length int, options ...func(setter optionsSetters)) Validator[string] {
	validator := &stringLengthValidator{length: length}

	for _, option := range options {
		option(validator)
	}

	return validator
}

type stringLengthValidator struct {
	length       int
	errorMessage string
}

func (validator *stringLengthValidator) SetErrorMessage(err string) {
	validator.errorMessage = err
}

func (validator *stringLengthValidator) Validate(value string) Errors {
	if utf8.RuneCountInString(value) != validator.length {
		errorMessage := validator.errorMessage
		if errorMessage == "" {
			errorMessage = fmt.Sprintf("must have exactly %d characters", validator.length)
		}
		return []string{errorMessage}
	}
	return nil
}

// MinRunes returns a validator that ensures the string contains at least
// `min` runes (Unicode code points).
func MinRunes(min int, options ...func(setter optionsSetters)) Validator[string] {
	validator := &stringMinValidator{min: min}

	for _, option := range options {
		option(validator)
	}

	return validator
}

type stringMinValidator struct {
	min          int
	errorMessage string
}

func (validator *stringMinValidator) SetErrorMessage(err string) {
	validator.errorMessage = err
}

func (validator *stringMinValidator) Validate(value string) Errors {
	if utf8.RuneCountInString(value) < validator.min {
		errorMessage := validator.errorMessage
		if errorMessage == "" {
			errorMessage = fmt.Sprintf("must have at least %d characters", validator.min)
		}
		return []string{errorMessage}
	}
	return nil
}

// MaxRunes returns a validator that ensures the string contains at most
// `max` runes (Unicode code points).
func MaxRunes(max int, options ...func(setter optionsSetters)) Validator[string] {
	validator := &stringMaxValidator{max: max}

	for _, option := range options {
		option(validator)
	}

	return validator
}

type stringMaxValidator struct {
	max          int
	errorMessage string
}

func (validator *stringMaxValidator) SetErrorMessage(err string) {
	validator.errorMessage = err
}

func (validator *stringMaxValidator) Validate(value string) Errors {
	if utf8.RuneCountInString(value) > validator.max {
		errorMessage := validator.errorMessage
		if errorMessage == "" {
			errorMessage = fmt.Sprintf("must have at most %d characters", validator.max)
		}
		return []string{errorMessage}
	}
	return nil
}
