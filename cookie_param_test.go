package ino

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCookieParamString(t *testing.T) {
	tests := []struct {
		name     string
		cookie   string
		value    string
		expected string
		hasError bool
	}{
		{
			name:     "valid string",
			cookie:   "test=value",
			value:    "value",
			expected: "value",
			hasError: false,
		},
		{
			name:     "empty string",
			cookie:   "test=",
			value:    "",
			expected: "",
			hasError: false,
		},
		{
			name:     "missing cookie",
			cookie:   "",
			value:    "",
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.cookie != "" {
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

func TestCookieParamBool(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
		hasError bool
	}{
		{
			name:     "true",
			value:    "true",
			expected: true,
			hasError: false,
		},
		{
			name:     "false",
			value:    "false",
			expected: false,
			hasError: false,
		},
		{
			name:     "1 as true",
			value:    "1",
			expected: true,
			hasError: false,
		},
		{
			name:     "0 as false",
			value:    "0",
			expected: false,
			hasError: false,
		},
		{
			name:     "invalid bool",
			value:    "maybe",
			expected: false,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: "test", Value: tt.value})

			param := CookieParam[bool]("test")
			val, errs := param.GetParamValue(req)

			if tt.hasError && len(errs) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}
			if val != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, val)
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
