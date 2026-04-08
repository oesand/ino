package mux_test

import (
	"net/http"
	"testing"

	"github.com/oesand/ino/mux"
)

func TestIsValidMethod(t *testing.T) {
	tests := []struct {
		method string
		valid  bool
	}{
		{http.MethodGet, true},
		{http.MethodPost, true},
		{http.MethodPut, true},
		{http.MethodDelete, true},
		{http.MethodOptions, true},
		{http.MethodHead, true},
		{http.MethodConnect, true},
		{http.MethodPatch, true},
		{http.MethodTrace, true},
		{"INVALID", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			result := mux.IsValidMethod(tt.method)
			if result != tt.valid {
				t.Errorf("expected %v, got %v for method %s", tt.valid, result, tt.method)
			}
		})
	}
}
