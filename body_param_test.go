package ino

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/oesand/ino/validate"
)

func TestRequestParam(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	param := RequestParam()
	result, errs := param.GetParamValue(req)

	if len(errs) > 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
	if result != req {
		t.Errorf("expected request object to be returned unchanged")
	}
}

func TestRequestParam_Optional(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	param := RequestParam().Optional()
	result, errs := param.GetParamValue(req)

	if len(errs) > 0 {
		t.Errorf("expected no errors for optional param, got: %v", errs)
	}
	if result != req {
		t.Errorf("expected request object to be returned unchanged")
	}
}

func TestBodyParam(t *testing.T) {
	tests := []struct {
		name     string
		body     io.Reader
		optional bool
		hasError bool
	}{
		{
			name:     "valid body",
			body:     strings.NewReader("test content"),
			optional: false,
			hasError: false,
		},
		{
			name:     "nil body required",
			body:     nil,
			optional: false,
			hasError: true,
		},
		{
			name:     "nil body optional",
			body:     nil,
			optional: true,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", tt.body)
			if tt.body == nil {
				req.Body = nil
			}

			param := BodyParam()
			if tt.optional {
				param = param.Optional()
			}

			result, errs := param.GetParamValue(req)

			if tt.hasError && len(errs) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}

			if !tt.hasError && tt.body != nil {
				if result == nil {
					t.Errorf("expected non-nil body reader")
				}
			}
		})
	}
}

func TestMultipartFormParam(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		boundary  string
		maxMemory int64
		optional  bool
		hasError  bool
	}{
		{
			name:      "valid multipart form",
			body:      "--boundary\r\nContent-Disposition: form-data; name=\"field\"\r\n\r\nvalue\r\n--boundary--",
			boundary:  "boundary",
			maxMemory: 1024,
			optional:  false,
			hasError:  false,
		},
		{
			name:      "invalid multipart form required",
			body:      "invalid multipart data",
			boundary:  "",
			maxMemory: 1024,
			optional:  false,
			hasError:  true,
		},
		{
			name:      "invalid multipart form optional",
			body:      "invalid multipart data",
			boundary:  "",
			maxMemory: 1024,
			optional:  true,
			hasError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", strings.NewReader(tt.body))
			if tt.boundary != "" {
				req.Header.Set("Content-Type", "multipart/form-data; boundary="+tt.boundary)
			}

			param := MultipartFormParam(tt.maxMemory)
			if tt.optional {
				param = param.Optional()
			}

			result, errs := param.GetParamValue(req)

			if tt.hasError && len(errs) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}

			if !tt.hasError && tt.boundary != "" {
				if result == nil {
					t.Errorf("expected non-nil multipart form")
				}
			}
		})
	}
}

func TestJsonParam(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name     string
		body     string
		optional bool
		expected *TestStruct
		hasError bool
	}{
		{
			name:     "valid json",
			body:     `{"name":"John","age":30}`,
			optional: false,
			expected: &TestStruct{Name: "John", Age: 30},
			hasError: false,
		},
		{
			name:     "invalid json required",
			body:     `{"name":"John","age":}`,
			optional: false,
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid json optional",
			body:     `{"name":"John","age":}`,
			optional: true,
			expected: nil,
			hasError: false,
		},
		{
			name:     "nil body required",
			body:     "",
			optional: false,
			expected: nil,
			hasError: true,
		},
		{
			name:     "nil body optional",
			body:     "",
			optional: true,
			expected: nil,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader
			if tt.body != "" {
				body = strings.NewReader(tt.body)
			}

			req := httptest.NewRequest("POST", "/test", body)
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			param := JsonParam[TestStruct]()
			if tt.optional {
				param = param.Optional()
			}

			result, errs := param.GetParamValue(req)

			if tt.hasError && len(errs) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}

			if tt.expected != nil && result != nil {
				if result.Name != tt.expected.Name || result.Age != tt.expected.Age {
					t.Errorf("expected %+v, got %+v", tt.expected, result)
				}
			} else if tt.expected != result {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestJsonParam_WithValidators(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// Validator that requires name to be non-empty
	nameValidator := validate.FuncValidator[*TestStruct](func(val *TestStruct) validate.Errors {
		if val.Name == "" {
			return []string{"name is required"}
		}
		return nil
	})

	req := httptest.NewRequest("POST", "/test", strings.NewReader(`{"name":"","age":25}`))
	req.Header.Set("Content-Type", "application/json")

	param := JsonParam[TestStruct](nameValidator)
	result, errs := param.GetParamValue(req)

	if len(errs) == 0 {
		t.Errorf("expected validation error for empty name")
	}
	if result != nil {
		t.Errorf("expected nil result when validation fails")
	}
	if !strings.Contains(strings.Join(errs, ""), "name is required") {
		t.Errorf("expected 'name is required' error, got: %v", errs)
	}
}
