package validate

type optionsSetters interface {
	SetErrorMessage(err string)
}

// WithMessage returns an option which sets a custom error message on a
// validator that supports message configuration.
func WithMessage(message string) func(optionsSetters) {
	return func(s optionsSetters) {
		s.SetErrorMessage(message)
	}
}
