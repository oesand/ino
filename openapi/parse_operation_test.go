package openapi

import (
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/oesand/ino"
)

func TestNewOperationUsesClearedPatternAndDefaults(t *testing.T) {
	route := ino.New(ino.Get("/users/{id:\\d+}", ino.F(func(http.ResponseWriter, *http.Request) {})))

	var first ino.Route
	for r := range route.Routes() {
		first = r
		break
	}

	pattern, operation := newOperation(first)
	if pattern != "/users/{id}" {
		t.Fatalf("expected cleared pattern /users/{id}, got %s", pattern)
	}
	if operation.ID != "get_users_id" {
		t.Fatalf("expected operation id get_users_id, got %s", operation.ID)
	}
	if operation.Summary != "GET /users/{id}" {
		t.Fatalf("expected default summary, got %s", operation.Summary)
	}
	if operation.Responses == nil || operation.Responses.StatusCodeResponses == nil {
		t.Fatal("expected initialized responses map")
	}
}

func TestOperationApply(t *testing.T) {
	tests := []struct {
		method string
		has    func(spec.PathItem) bool
	}{
		{http.MethodGet, func(path spec.PathItem) bool { return path.Get != nil }},
		{http.MethodPost, func(path spec.PathItem) bool { return path.Post != nil }},
		{http.MethodPut, func(path spec.PathItem) bool { return path.Put != nil }},
		{http.MethodDelete, func(path spec.PathItem) bool { return path.Delete != nil }},
		{http.MethodPatch, func(path spec.PathItem) bool { return path.Patch != nil }},
		{http.MethodHead, func(path spec.PathItem) bool { return path.Head != nil }},
		{http.MethodOptions, func(path spec.PathItem) bool { return path.Options != nil }},
		{"TRACE", func(path spec.PathItem) bool { return path.Get != nil }},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			operation := operationWrapper(spec.Operation{})
			var path spec.PathItem
			operation.Apply(&path, tt.method)
			if !tt.has(path) {
				t.Fatalf("expected method %s to be applied to path item", tt.method)
			}
		})
	}
}

func TestOperationFillConsumes(t *testing.T) {
	wrapper := newSwagger(&spec.Info{})
	mux := ino.New(ino.Post("/users/{id}/avatar", ino.ParamHandler1(
		ino.FileParam("avatar"),
		func(_ *multipart.FileHeader, _ http.ResponseWriter) {},
	)))

	var route ino.Route
	for r := range mux.Routes() {
		route = r
		break
	}

	operation := operationWrapper(spec.Operation{})
	operation.FillConsumes(route, wrapper)

	if len(operation.Parameters) != 2 {
		t.Fatalf("expected file and path parameters, got %#v", operation.Parameters)
	}
	params := map[string]spec.Parameter{}
	for _, param := range operation.Parameters {
		params[param.Name] = param
	}
	if params["avatar"].In != "formData" || params["avatar"].Type != "file" {
		t.Fatalf("expected avatar file parameter, got %#v", params["avatar"])
	}
	if params["id"].In != "path" || params["id"].Type != "string" {
		t.Fatalf("expected inferred path parameter, got %#v", params["id"])
	}
	if len(operation.Consumes) != 1 || operation.Consumes[0] != "multipart/form-data" {
		t.Fatalf("expected multipart consume, got %#v", operation.Consumes)
	}
}

func TestOperationScanAttributes(t *testing.T) {
	wrapper := newSwagger(&spec.Info{})
	parameter := Header[string]("X-Request-ID")
	route := ino.Get("/users", ino.F(func(http.ResponseWriter, *http.Request) {}),
		Respond[attrsTestResponse](200).WithDescription("ok"),
		RespondHtml(404),
		Tag("users"),
		Summary("list users"),
		Auth("oauth2", "users:read"),
		parameter,
	)

	operation := operationWrapper(spec.Operation{
		OperationProps: spec.OperationProps{
			Responses: &spec.Responses{ResponsesProps: spec.ResponsesProps{StatusCodeResponses: map[int]spec.Response{}}},
		},
	})
	operation.ScanAttributes(route, wrapper)

	if operation.Summary != "list users" {
		t.Fatalf("expected custom summary, got %s", operation.Summary)
	}
	if len(operation.Tags) != 1 || operation.Tags[0] != "users" {
		t.Fatalf("expected users tag, got %#v", operation.Tags)
	}
	if len(operation.Security) != 1 || len(operation.Security[0]["oauth2"]) != 1 {
		t.Fatalf("expected oauth2 security, got %#v", operation.Security)
	}
	if len(operation.Parameters) != 1 || operation.Parameters[0].Name != "X-Request-ID" {
		t.Fatalf("expected custom header parameter, got %#v", operation.Parameters)
	}
	if operation.Responses.StatusCodeResponses[200].Description != "ok" {
		t.Fatalf("expected ok description, got %#v", operation.Responses.StatusCodeResponses[200])
	}
	if operation.Responses.StatusCodeResponses[200].Schema.Ref.String() == "" {
		t.Fatalf("expected response schema ref, got %#v", operation.Responses.StatusCodeResponses[200].Schema)
	}
	if operation.Responses.StatusCodeResponses[404].Schema.Type[0] != "string" {
		t.Fatalf("expected html response string schema, got %#v", operation.Responses.StatusCodeResponses[404].Schema)
	}

	produces := map[string]bool{}
	for _, value := range operation.Produces {
		produces[value] = true
	}
	if !produces["application/json"] || !produces["text/html"] {
		t.Fatalf("expected json and html produces, got %#v", operation.Produces)
	}
}

func TestGenerateOperationID(t *testing.T) {
	tests := map[string]string{
		"/":                         "post_root",
		"/users/{id}":               "post_users_id",
		"/users/{id}/profile-image": "post_users_id_profile_image",
	}

	for path, expected := range tests {
		if got := generateOperationID(http.MethodPost, path); got != expected {
			t.Fatalf("expected %s for %s, got %s", expected, path, got)
		}
	}
}
