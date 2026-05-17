package validate

// Rule creates a validator from a predicate function.
//
// The returned Validator succeeds when the provided rule returns true and
// fails when it returns false. A custom error message can be provided using
// WithMessage.
func Rule[Value any](rule func(Value) bool, options ...func(setter optionsSetters)) Validator[Value] {
	validator := &ruleValidator[Value]{rule: rule}

	for _, option := range options {
		option(validator)
	}

	return validator
}

type ruleValidator[Value any] struct {
	rule         func(Value) bool
	errorMessage string
}

func (validator *ruleValidator[Value]) SetErrorMessage(err string) {
	validator.errorMessage = err
}

func (validator *ruleValidator[Value]) Validate(value Value) Errors {
	if !validator.rule(value) {
		errorMessage := validator.errorMessage
		if errorMessage == "" {
			errorMessage = "defined rule was not met"
		}
		return []string{errorMessage}
	}
	return nil
}
