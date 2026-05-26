package ino

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/oesand/ino/validate"
)

func TestCookieParamString(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
		hasError bool
	}{
		{
			name:     "valid string",
			value:    "value",
			expected: "value",
			hasError: false,
		},
		{
			name:     "empty string",
			value:    "",
			expected: "",
			hasError: true,
		},
		{
			name:     "missing cookie",
			value:    "",
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.value != "" {
				req.AddCookie(&http.Cookie{Name: "test", Value: tt.value})
			}

			param := CookieParam[string]("test")
			val, errs := param.GetParamValue(req)

			if tt.hasError && len(errs) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}
			if val != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, val)
			}
		})
	}
}

func TestCookieParamInt64(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected int64
		hasError bool
	}{
		{
			name:     "valid int64",
			value:    "42",
			expected: 42,
			hasError: false,
		},
		{
			name:     "zero",
			value:    "0",
			expected: 0,
			hasError: false,
		},
		{
			name:     "negative",
			value:    "-10",
			expected: -10,
			hasError: false,
		},
		{
			name:     "invalid int64",
			value:    "not-a-number",
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: "test", Value: tt.value})

			param := CookieParam[int64]("test")
			val, errs := param.GetParamValue(req)

			if tt.hasError && len(errs) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}
			if val != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, val)
			}
		})
	}
}

func TestCookieParam_Optional(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	// No cookie set

	param := CookieParam[string]("test").Optional()
	val, errs := param.GetParamValue(req)

	if len(errs) > 0 {
		t.Errorf("expected no error for optional param but got: %v", errs)
	}
	if val != "" {
		t.Errorf("expected empty string for missing optional param, got %q", val)
	}
}

// TestCookieParamWithRuleValidator tests CookieParam with a custom rule validator
func TestCookieParamWithRuleValidator(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		hasError bool
	}{
		{
			name:     "value meets rule",
			value:    "hello",
			hasError: false,
		},
		{
			name:     "value fails rule",
			value:    "hi",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: "name", Value: tt.value})

			validator := validate.Rule[string](func(s string) bool {
				return len(s) > 2
			}, validate.WithMessage("must be longer than 2 characters"))

			param := CookieParam[string]("name", validator)
			val, errs := param.GetParamValue(req)

			if tt.hasError && len(errs) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}
			if val != tt.value {
				t.Errorf("expected %q, got %q", tt.value, val)
			}
		})
	}
}

// TestCookieParamWithMinValidator tests CookieParam with Min validator for int64
func TestCookieParamWithMinValidator(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		hasError bool
	}{
		{
			name:     "value >= min",
			value:    "100",
			hasError: false,
		},
		{
			name:     "value equals min",
			value:    "50",
			hasError: false,
		},
		{
			name:     "value < min",
			value:    "25",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: "age", Value: tt.value})

			param := CookieParam[int64]("age", validate.Min[int64](50))
			_, errs := param.GetParamValue(req)

			if tt.hasError && len(errs) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}
		})
	}
}

// TestCookieParamWithMaxValidator tests CookieParam with Max validator for int64
func TestCookieParamWithMaxValidator(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		hasError bool
	}{
		{
			name:     "value <= max",
			value:    "50",
			hasError: false,
		},
		{
			name:     "value equals max",
			value:    "100",
			hasError: false,
		},
		{
			name:     "value > max",
			value:    "150",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: "port", Value: tt.value})

			param := CookieParam[int64]("port", validate.Max[int64](100))
			_, errs := param.GetParamValue(req)

			if tt.hasError && len(errs) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}
		})
	}
}

// TestCookieParamWithMultipleValidators tests CookieParam with multiple validators
func TestCookieParamWithMultipleValidators(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		hasError bool
	}{
		{
			name:     "passes all validators",
			value:    "good_name",
			hasError: false,
		},
		{
			name:     "fails min length",
			value:    "ab",
			hasError: true,
		},
		{
			name:     "fails regex pattern",
			value:    "bad-name",
			hasError: true,
		},
		{
			name:     "fails multiple validators",
			value:    "x",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: "username", Value: tt.value})

			minLengthValidator := validate.Rule[string](
				func(s string) bool { return len(s) >= 3 },
				validate.WithMessage("must be at least 3 characters"),
			)
			regexValidator := validate.Regex(
				regexp.MustCompile("^[a-z_]+$"),
				validate.WithMessage("must contain only lowercase letters and underscores"),
			)

			param := CookieParam[string]("username", minLengthValidator, regexValidator)
			val, errs := param.GetParamValue(req)

			if tt.hasError && len(errs) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}
			if val != tt.value {
				t.Errorf("expected %q, got %q", tt.value, val)
			}
		})
	}
}

// TestCookieParamOptionalWithValidator tests optional parameters with validators
func TestCookieParamOptionalWithValidator(t *testing.T) {
	tests := []struct {
		name       string
		setCookie  bool
		value      string
		hasError   bool
		expectZero bool
	}{
		{
			name:       "optional param not provided",
			setCookie:  false,
			value:      "",
			hasError:   false,
			expectZero: true,
		},
		{
			name:       "optional param provided and valid",
			setCookie:  true,
			value:      "75",
			hasError:   false,
			expectZero: false,
		},
		{
			name:       "optional param provided but fails validation",
			setCookie:  true,
			value:      "25",
			hasError:   true,
			expectZero: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.setCookie {
				req.AddCookie(&http.Cookie{Name: "priority", Value: tt.value})
			}

			validator := validate.Min[int64](50)
			param := CookieParam[int64]("priority", validator).Optional()
			val, errs := param.GetParamValue(req)

			if tt.hasError && len(errs) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}
			if tt.expectZero && val != 0 {
				t.Errorf("expected 0 for missing optional param, got %d", val)
			}
		})
	}
}

// TestCookieParamBoolVariations tests various boolean string representations
func TestCookieParamBoolVariations(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
		hasError bool
	}{
		{
			name:     "true lowercase",
			value:    "true",
			expected: true,
			hasError: false,
		},
		{
			name:     "false lowercase",
			value:    "false",
			expected: false,
			hasError: false,
		},
		{
			name:     "TRUE uppercase",
			value:    "TRUE",
			expected: true,
			hasError: false,
		},
		{
			name:     "FALSE uppercase",
			value:    "FALSE",
			expected: false,
			hasError: false,
		},
		{
			name:     "t short form",
			value:    "t",
			expected: true,
			hasError: false,
		},
		{
			name:     "f short form",
			value:    "f",
			expected: false,
			hasError: false,
		},
		{
			name:     "yes string",
			value:    "yes",
			expected: false,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: "flag", Value: tt.value})

			param := CookieParam[bool]("flag")
			val, errs := param.GetParamValue(req)

			if tt.hasError && len(errs) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}
			if !tt.hasError && val != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, val)
			}
		})
	}
}

// TestCookieParamIntParsing tests various integer parsing scenarios
func TestCookieParamIntParsing(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected int64
		hasError bool
	}{
		{
			name:     "leading zeros",
			value:    "00042",
			expected: 42,
			hasError: false,
		},
		{
			name:     "whitespace around value",
			value:    "  42  ",
			expected: 0,
			hasError: true,
		},
		{
			name:     "plus sign",
			value:    "+42",
			expected: 42,
			hasError: false,
		},
		{
			name:     "decimal point",
			value:    "42.5",
			expected: 0,
			hasError: true,
		},
		{
			name:     "hex format",
			value:    "0x2a",
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: "number", Value: tt.value})

			param := CookieParam[int64]("number")
			val, errs := param.GetParamValue(req)

			if tt.hasError && len(errs) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}
			if !tt.hasError && val != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, val)
			}
		})
	}
}
