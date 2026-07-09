package openapi

import (
	"net/http"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/oesand/mo"
)

type generateSchemaUser struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func TestGenerateSchema(t *testing.T) {
	mux := mo.New(
		mo.Get("/users/{id:\\d+}", mo.ParamHandler1(
			mo.PathParam[int64]("id"),
			func(int64, http.ResponseWriter) {},
		), Respond[generateSchemaUser](200), Tag("users"), Summary("get user")),
		mo.Post("/users", mo.ParamHandler1(
			mo.JsonParam[generateSchemaUser](),
			func(*generateSchemaUser, http.ResponseWriter) {},
		), Respond[generateSchemaUser](201).WithDescription("created")),
	)

	schema, err := GenerateSchema(mux, &GenerateOptions{
		Info: spec.InfoProps{Title: "Users API", Version: "2.0.0"},
		Tags: []spec.Tag{{TagProps: spec.TagProps{Name: "users"}}},
	})
	if err != nil {
		t.Fatalf("unexpected generate error: %v", err)
	}

	if schema.Swagger != "2.0" {
		t.Fatalf("expected swagger 2.0, got %s", schema.Swagger)
	}
	if schema.Info.Title != "Users API" || schema.Info.Version != "2.0.0" {
		t.Fatalf("unexpected info: %#v", schema.Info)
	}
	if len(schema.Tags) != 1 || schema.Tags[0].Name != "users" {
		t.Fatalf("expected users tag, got %#v", schema.Tags)
	}

	getOperation := schema.Paths.Paths["/users/{id}"].Get
	if getOperation == nil {
		t.Fatal("expected GET /users/{id} operation")
	}
	if getOperation.Summary != "get user" {
		t.Fatalf("expected custom summary, got %s", getOperation.Summary)
	}
	if len(getOperation.Parameters) != 1 || getOperation.Parameters[0].Name != "id" {
		t.Fatalf("expected id path parameter, got %#v", getOperation.Parameters)
	}
	if getOperation.Parameters[0].Type != "integer" || getOperation.Parameters[0].Format != "int64" {
		t.Fatalf("expected int64 path parameter, got %#v", getOperation.Parameters[0])
	}
	if getOperation.Responses.StatusCodeResponses[200].Schema.Ref.String() == "" {
		t.Fatalf("expected response ref schema, got %#v", getOperation.Responses.StatusCodeResponses[200].Schema)
	}

	postOperation := schema.Paths.Paths["/users"].Post
	if postOperation == nil {
		t.Fatal("expected POST /users operation")
	}
	if len(postOperation.Consumes) != 1 || postOperation.Consumes[0] != "application/json" {
		t.Fatalf("expected json consume, got %#v", postOperation.Consumes)
	}
	if postOperation.Responses.StatusCodeResponses[201].Description != "created" {
		t.Fatalf("expected created response description, got %#v", postOperation.Responses.StatusCodeResponses[201])
	}
	if len(schema.Definitions) == 0 {
		t.Fatal("expected generated definitions")
	}
}

func TestGenerateSchemaDefaultsAndErrors(t *testing.T) {
	if _, err := GenerateSchema(nil, nil); err == nil {
		t.Fatal("expected nil mux error")
	}

	schema, err := GenerateSchema(mo.New(), nil)
	if err != nil {
		t.Fatalf("unexpected generate error: %v", err)
	}
	if schema.Info.Title != "API" {
		t.Fatalf("expected default title API, got %s", schema.Info.Title)
	}
	if schema.Info.Version != "1.0.0" {
		t.Fatalf("expected default version 1.0.0, got %s", schema.Info.Version)
	}
}
