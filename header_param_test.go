package ino_test

import (
	"net/http"
	"testing"

	"github.com/oesand/ino"
)

// TestHeaderParamString tests HeaderParam with string type.
func TestHeaderParamString(t *testing.T) {
	tests := []struct {
		name        string
		headerName  string
		headerValue string
		expectError bool
		errorMsg    string
		isOptional  bool
	}{
		{
			name:        "valid string header",
			headerName:  "X-Test",
			headerValue: "application/json",
		},
		{
			name:        "missing required string header",
			headerName:  "Authorization",
			expectError: true,
			errorMsg:    "header 'Authorization' is required",
		},
		{
			name:       "optional missing string header",
			headerName: "Optional",
			isOptional: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.headerValue != "" {
				req.Header.Set(tt.headerName, tt.headerValue)
			}

			provider := ino.HeaderParam[string](tt.headerName)
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
				if !tt.isOptional && tt.headerValue != "" && val != tt.headerValue {
					t.Errorf("expected %v, got %v", tt.headerValue, val)
				}
			}
		})
	}
}

// TestHeaderParamInt64 tests HeaderParam with int64 type.
func TestHeaderParamInt64(t *testing.T) {
	tests := []struct {
		name        string
		headerName  string
		headerValue string
		expected    int64
		expectError bool
		errorMsg    string
		isOptional  bool
	}{
		{
			name:        "valid int64 header",
			headerName:  "Content-Length",
			headerValue: "1024",
			expected:    1024,
		},
		{
			name:        "invalid int64 header",
			headerName:  "Content-Length",
			headerValue: "notint",
			expectError: true,
			errorMsg:    "header 'Content-Length' must be integer",
		},
		{
			name:        "missing required int64 header",
			headerName:  "X-Count",
			expectError: true,
			errorMsg:    "header 'X-Count' is required",
		},
		{
			name:       "optional missing int64 header",
			headerName: "X-Optional",
			isOptional: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.headerValue != "" {
				req.Header.Set(tt.headerName, tt.headerValue)
			}

			provider := ino.HeaderParam[int64](tt.headerName)
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
				if !tt.isOptional && tt.headerValue != "" && val != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, val)
				}
			}
		})
	}
}

// TestHeaderParamBool tests HeaderParam with bool type.
func TestHeaderParamBool(t *testing.T) {
	tests := []struct {
		name        string
		headerName  string
		headerValue string
		expected    bool
		expectError bool
		errorMsg    string
		isOptional  bool
	}{
		{
			name:        "bool header true",
			headerName:  "X-Debug",
			headerValue: "true",
			expected:    true,
		},
		{
			name:        "bool header false",
			headerName:  "X-Debug",
			headerValue: "false",
			expected:    false,
		},
		{
			name:        "invalid bool header",
			headerName:  "X-Debug",
			headerValue: "notbool",
			expectError: true,
			errorMsg:    "header 'X-Debug' must be bool",
		},
		{
			name:        "missing required bool header",
			headerName:  "X-Enabled",
			expectError: true,
			errorMsg:    "header 'X-Enabled' is required",
		},
		{
			name:       "optional missing bool header",
			headerName: "X-Optional",
			isOptional: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.headerValue != "" {
				req.Header.Set(tt.headerName, tt.headerValue)
			}

			provider := ino.HeaderParam[bool](tt.headerName)
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
				if !tt.isOptional && tt.headerValue != "" && val != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, val)
				}
			}
		})
	}
}

// TestHeaderParam_Optional tests the Optional method specifically.
func TestHeaderParam_Optional(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	provider := ino.HeaderParam[string]("missing").Optional()
	val, errs := provider.GetParamValue(req)

	if len(errs) > 0 {
		t.Errorf("optional header should not error, got %v", errs)
	}
	if val != "" {
		t.Errorf("expected empty string for missing optional header, got %q", val)
	}
}
