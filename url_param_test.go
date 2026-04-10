package ino_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/oesand/ino"
	"github.com/oesand/ino/internal"
)

// TestUrlParamString tests UrlParam with string type.
func TestUrlParamString(t *testing.T) {
	tests := []struct {
		name        string
		paramName   string
		paramValue  string
		expectError bool
		errorMsg    string
		isOptional  bool
	}{
		{
			name:       "valid string param",
			paramName:  "name",
			paramValue: "john",
		},
		{
			name:        "missing required string param",
			paramName:   "missing",
			expectError: true,
			errorMsg:    "url param 'missing' is required",
		},
		{
			name:       "optional missing string param",
			paramName:  "optional",
			isOptional: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			urlParams := make(map[string]string)
			if tt.paramValue != "" {
				urlParams[tt.paramName] = tt.paramValue
			}
			ctx = context.WithValue(ctx, internal.CtxKey{Key: "mux/url_params"}, urlParams)

			req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)

			provider := ino.UrlParam[string](tt.paramName)
			if tt.isOptional {
				provider = provider.Optional()
			}
			val, errs := provider.GetParamValue(req)

			if tt.expectError {
				if len(errs) == 0 {
					t.Errorf("expected error, got none")
				} else if errs[0] != tt.errorMsg {
					t.Errorf("expected error %q, got %q", tt.errorMsg, errs[0])
				}
			} else {
				if len(errs) > 0 {
					t.Errorf("unexpected error: %v", errs)
				}
				if !tt.isOptional && tt.paramValue != "" && val != tt.paramValue {
					t.Errorf("expected %v, got %v", tt.paramValue, val)
				}
			}
		})
	}
}

// TestUrlParamInt64 tests UrlParam with int64 type.
func TestUrlParamInt64(t *testing.T) {
	tests := []struct {
		name        string
		paramName   string
		paramValue  string
		expected    int64
		expectError bool
		errorMsg    string
		isOptional  bool
	}{
		{
			name:       "valid int64 param",
			paramName:  "age",
			paramValue: "25",
			expected:   25,
		},
		{
			name:        "invalid int64 param",
			paramName:   "age",
			paramValue:  "notint",
			expectError: true,
			errorMsg:    "url param 'age' must be integer",
		},
		{
			name:        "missing required int64 param",
			paramName:   "missing",
			expectError: true,
			errorMsg:    "url param 'missing' is required",
		},
		{
			name:       "optional missing int64 param",
			paramName:  "optional",
			isOptional: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			urlParams := make(map[string]string)
			if tt.paramValue != "" {
				urlParams[tt.paramName] = tt.paramValue
			}
			ctx = context.WithValue(ctx, internal.CtxKey{Key: "mux/url_params"}, urlParams)

			req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)

			provider := ino.UrlParam[int64](tt.paramName)
			if tt.isOptional {
				provider = provider.Optional()
			}
			val, errs := provider.GetParamValue(req)

			if tt.expectError {
				if len(errs) == 0 {
					t.Errorf("expected error, got none")
				} else if errs[0] != tt.errorMsg {
					t.Errorf("expected error %q, got %q", tt.errorMsg, errs[0])
				}
			} else {
				if len(errs) > 0 {
					t.Errorf("unexpected error: %v", errs)
				}
				if !tt.isOptional && tt.paramValue != "" && val != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, val)
				}
			}
		})
	}
}

// TestUrlParamBool tests UrlParam with bool type.
func TestUrlParamBool(t *testing.T) {
	tests := []struct {
		name        string
		paramName   string
		paramValue  string
		expected    bool
		expectError bool
		errorMsg    string
		isOptional  bool
	}{
		{
			name:       "bool param true",
			paramName:  "active",
			paramValue: "true",
			expected:   true,
		},
		{
			name:       "bool param false",
			paramName:  "active",
			paramValue: "false",
			expected:   false,
		},
		{
			name:        "invalid bool param",
			paramName:   "active",
			paramValue:  "notbool",
			expectError: true,
			errorMsg:    "url param 'active' must be bool",
		},
		{
			name:        "missing required bool param",
			paramName:   "missing",
			expectError: true,
			errorMsg:    "url param 'missing' is required",
		},
		{
			name:       "optional missing bool param",
			paramName:  "optional",
			isOptional: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			urlParams := make(map[string]string)
			if tt.paramValue != "" {
				urlParams[tt.paramName] = tt.paramValue
			}
			ctx = context.WithValue(ctx, internal.CtxKey{Key: "mux/url_params"}, urlParams)

			req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)

			provider := ino.UrlParam[bool](tt.paramName)
			if tt.isOptional {
				provider = provider.Optional()
			}
			val, errs := provider.GetParamValue(req)

			if tt.expectError {
				if len(errs) == 0 {
					t.Errorf("expected error, got none")
				} else if errs[0] != tt.errorMsg {
					t.Errorf("expected error %q, got %q", tt.errorMsg, errs[0])
				}
			} else {
				if len(errs) > 0 {
					t.Errorf("unexpected error: %v", errs)
				}
				if !tt.isOptional && tt.paramValue != "" && val != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, val)
				}
			}
		})
	}
}

// TestUrlParam_Optional tests the Optional method specifically.
func TestUrlParam_Optional(t *testing.T) {
	ctx := context.Background()
	urlParams := make(map[string]string)
	ctx = context.WithValue(ctx, internal.CtxKey{Key: "mux/url_params"}, urlParams)

	req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)

	provider := ino.UrlParam[string]("missing").Optional()
	val, errs := provider.GetParamValue(req)

	if len(errs) > 0 {
		t.Errorf("optional param should not error, got %v", errs)
	}
	if val != "" {
		t.Errorf("expected empty string for missing optional param, got %v", val)
	}
}
