package openapi

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/oesand/mo"
	"github.com/oesand/mo/internal"
)

type swaggerTestBody struct {
	Name string `json:"name"`
}

func TestNewSwagger(t *testing.T) {
	info := &spec.Info{InfoProps: spec.InfoProps{Title: "Test API"}}
	swagger := newSwagger(info)

	if swagger.Swagger != "2.0" {
		t.Fatalf("expected swagger 2.0, got %s", swagger.Swagger)
	}
	if swagger.Info != info {
		t.Fatal("expected info pointer to be preserved")
	}
	if swagger.Paths == nil || swagger.Paths.Paths == nil {
		t.Fatal("expected initialized paths")
	}
	if swagger.Definitions == nil {
		t.Fatal("expected initialized definitions")
	}
	if swagger.Base() == nil {
		t.Fatal("expected base swagger")
	}
}

func TestMakeSchemaRegistersNamedStructDefinitions(t *testing.T) {
	swagger := newSwagger(&spec.Info{})

	schema := swagger.MakeSchema(reflect.TypeFor[*swaggerTestBody]())
	if schema.Ref.String() == "" {
		t.Fatalf("expected ref schema, got %#v", schema)
	}
	if len(swagger.Definitions) != 1 {
		t.Fatalf("expected one definition, got %#v", swagger.Definitions)
	}

	sameSchema := swagger.MakeSchema(reflect.TypeFor[swaggerTestBody]())
	if sameSchema.Ref.String() != schema.Ref.String() {
		t.Fatalf("expected stable ref, got %s and %s", schema.Ref.String(), sameSchema.Ref.String())
	}
	if len(swagger.Definitions) != 1 {
		t.Fatalf("expected definition to be reused, got %#v", swagger.Definitions)
	}

	if got := swagger.MakeSchema(nil); got != nil {
		t.Fatalf("expected nil schema for nil type, got %#v", got)
	}
}

func TestMakeParamFromHandlerParam(t *testing.T) {
	swagger := newSwagger(&spec.Info{})

	tests := []struct {
		name       string
		param      any
		isUpdate   bool
		in         string
		paramType  string
		format     string
		consume    string
		required   bool
		collection string
	}{
		{"path", mo.PathParam[int64]("id"), false, "path", "integer", "int64", "", true, ""},
		{"header", mo.HeaderParam[bool]("X-Flag").Optional(), false, "header", "boolean", "", "", false, ""},
		{"cookie", mo.CookieParam[string]("session"), false, "cookie", "string", "", "", true, ""},
		{"form as query", mo.FormParam[int32]("page"), false, "query", "integer", "int32", "", true, ""},
		{"form as post form", mo.FormParam[string]("name"), true, "formData", "string", "", "application/x-www-form-urlencoded", true, ""},
		{"post form", mo.PostFormParam[string]("name"), false, "formData", "string", "", "application/x-www-form-urlencoded", true, ""},
		{"file", mo.FileParam("file"), false, "formData", "file", "", "multipart/form-data", false, ""},
		{"json", mo.JsonParam[swaggerTestBody](), false, "body", "", "", "application/json", true, ""},
		{"raw body", mo.BodyParam(), false, "body", "", "", "application/octet-stream", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param, consume := swagger.MakeParamFromHandlerParam(tt.param.(internal.ParamSchema), tt.isUpdate)
			if consume != tt.consume {
				t.Fatalf("expected consume %s, got %s", tt.consume, consume)
			}
			if param == nil {
				t.Fatal("expected parameter")
			}
			if param.In != tt.in {
				t.Fatalf("expected parameter in %s, got %s", tt.in, param.In)
			}
			if param.Type != tt.paramType {
				t.Fatalf("expected parameter type %s, got %s", tt.paramType, param.Type)
			}
			if param.Format != tt.format {
				t.Fatalf("expected format %s, got %s", tt.format, param.Format)
			}
			if param.Required != tt.required {
				t.Fatalf("expected required %v, got %v", tt.required, param.Required)
			}
			if tt.collection != "" && param.CollectionFormat != tt.collection {
				t.Fatalf("expected collection format %s, got %s", tt.collection, param.CollectionFormat)
			}
		})
	}
}

func TestMakeParamFromHandlerParamSpecialCases(t *testing.T) {
	swagger := newSwagger(&spec.Info{})

	route := mo.Post("/upload", mo.ParamHandler1(
		mo.MultipartFormParam(1024),
		func(*multipart.Form, http.ResponseWriter) {},
	))
	paramHandler := route.Handler().(interface{ Params() []internal.ParamSchema })
	param, consume := swagger.MakeParamFromHandlerParam(paramHandler.Params()[0], true)
	if param != nil {
		t.Fatalf("expected multipart form to add consume without parameter, got %#v", param)
	}
	if consume != "multipart/form-data" {
		t.Fatalf("expected multipart consume, got %s", consume)
	}
}

func TestMakeParamFromHandlerParamWithRoutes(t *testing.T) {
	swagger := newSwagger(&spec.Info{})
	route := mo.Post("/users/{id}", mo.ParamHandler1(
		mo.JsonParam[swaggerTestBody](),
		func(*swaggerTestBody, http.ResponseWriter) {},
	))

	paramHandler := route.Handler().(interface {
		Params() []internal.ParamSchema
	})
	params := paramHandler.Params()
	if len(params) != 1 {
		t.Fatalf("expected one handler param, got %d", len(params))
	}

	param, consume := swagger.MakeParamFromHandlerParam(params[0], true)
	if param.In != "body" || param.Name != "json" {
		t.Fatalf("expected json body parameter, got %#v", param)
	}
	if consume != "application/json" {
		t.Fatalf("expected application/json consume, got %s", consume)
	}
	if param.Schema.Ref.String() == "" {
		t.Fatalf("expected json body ref schema, got %#v", param.Schema)
	}
}

func TestRegister(t *testing.T) {
	mux := mo.New()

	err := Register(mux, nil, nil)
	if err == nil {
		t.Fatal("expected nil schema error")
	}

	schema := newSwagger(&spec.Info{InfoProps: spec.InfoProps{Title: "Docs", Version: "1.0.0"}}).Base()
	err = Register(mux, schema, &RegisterOptions{OpenApiPattern: "/docs/openapi.json", SwaggerPattern: "/docs/swagger/{*}"})
	if err != nil {
		t.Fatalf("unexpected register error: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/docs/openapi.json", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content type, got %s", got)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if len(body) == 0 {
		t.Fatal("expected openapi response body")
	}
}
