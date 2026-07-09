package openapi

import (
	"reflect"
	"testing"

	"github.com/go-openapi/spec"
)

type attrsTestResponse struct {
	ID int `json:"id"`
}

func TestRespond(t *testing.T) {
	attr := Respond[attrsTestResponse](201).
		WithContentType("application/vnd.api+json").
		WithDescription("created")

	if attr.Code != 201 {
		t.Fatalf("expected code 201, got %d", attr.Code)
	}
	if attr.ContentType != "application/vnd.api+json" {
		t.Fatalf("expected custom content type, got %s", attr.ContentType)
	}
	if attr.Description != "created" {
		t.Fatalf("expected description created, got %s", attr.Description)
	}
	if attr.Type != reflect.TypeFor[attrsTestResponse]() {
		t.Fatalf("expected response type %v, got %v", reflect.TypeFor[attrsTestResponse](), attr.Type)
	}
}

func TestRespondFile(t *testing.T) {
	attr := RespondFile(200)

	if attr.Code != 200 {
		t.Fatalf("expected code 200, got %d", attr.Code)
	}
	if attr.ContentType != "application/octet-stream" {
		t.Fatalf("expected octet-stream content type, got %s", attr.ContentType)
	}
	if got := attr.Schema.Type[0]; got != "file" {
		t.Fatalf("expected file schema, got %s", got)
	}
}

func TestRespondHtml(t *testing.T) {
	attr := RespondHtml(200)

	if attr.ContentType != "text/html" {
		t.Fatalf("expected text/html content type, got %s", attr.ContentType)
	}
	if got := attr.Schema.Type[0]; got != "string" {
		t.Fatalf("expected string schema, got %s", got)
	}
}

func TestRouteAttrs(t *testing.T) {
	tag := Tag("users")
	if tag.Name != "users" {
		t.Fatalf("expected tag users, got %s", tag.Name)
	}

	summary := Summary("list users")
	if summary.Summary != "list users" {
		t.Fatalf("expected summary list users, got %s", summary.Summary)
	}

	auth := Auth("oauth2", "users:read", "users:write")
	if auth.Name != "oauth2" {
		t.Fatalf("expected auth name oauth2, got %s", auth.Name)
	}
	if len(auth.Scopes) != 2 || auth.Scopes[0] != "users:read" || auth.Scopes[1] != "users:write" {
		t.Fatalf("unexpected auth scopes: %#v", auth.Scopes)
	}
}

func TestHeader(t *testing.T) {
	param := Header[int64]("X-Request-ID")

	if param.Name != "X-Request-ID" {
		t.Fatalf("expected header name X-Request-ID, got %s", param.Name)
	}
	if param.In != "header" {
		t.Fatalf("expected header parameter, got %s", param.In)
	}
	if param.Schema.Type[0] != "integer" || param.Schema.Format != "int64" {
		t.Fatalf("expected int64 schema, got type=%v format=%s", param.Schema.Type, param.Schema.Format)
	}
}

func TestCookie(t *testing.T) {
	param := Cookie[bool]("session")

	if param.Name != "session" {
		t.Fatalf("expected cookie name session, got %s", param.Name)
	}
	if param.In != "cookie" {
		t.Fatalf("expected cookie parameter, got %s", param.In)
	}
	if param.Schema.Type[0] != "boolean" {
		t.Fatalf("expected boolean schema, got %v", param.Schema.Type)
	}
}

func TestResponseAttrUsesCustomSchema(t *testing.T) {
	schema := spec.StringProperty()
	attr := &ResponseAttr{Code: 202, Schema: schema}

	if attr.Schema != schema {
		t.Fatal("expected response attr to keep custom schema")
	}
}
