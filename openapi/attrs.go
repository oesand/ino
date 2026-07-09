package openapi

import (
	"reflect"

	"github.com/go-openapi/spec"
)

// Respond declares a JSON response with a schema inferred from T.
func Respond[T any](code int) *ResponseAttr {
	return &ResponseAttr{
		Code:        code,
		ContentType: "application/json",
		Type:        reflect.TypeFor[T](),
	}
}

// RespondFile declares a binary file response.
func RespondFile(code int) *ResponseAttr {
	return &ResponseAttr{
		Code:        code,
		ContentType: "application/octet-stream",
		Schema:      new(spec.Schema).Typed("file", ""),
	}
}

// RespondHtml declares an HTML string response.
func RespondHtml(code int) *ResponseAttr {
	return &ResponseAttr{
		Code:        code,
		ContentType: "text/html",
		Schema:      new(spec.Schema).Typed("string", ""),
	}
}

// ResponseAttr describes an OpenAPI response attached to a route.
type ResponseAttr struct {
	Code        int
	ContentType string
	Type        reflect.Type
	Schema      *spec.Schema
	Description string
}

// WithContentType sets the response media type and returns the same attribute.
func (attr *ResponseAttr) WithContentType(value string) *ResponseAttr {
	attr.ContentType = value
	return attr
}

// WithDescription sets the response description and returns the same attribute.
func (attr *ResponseAttr) WithDescription(value string) *ResponseAttr {
	attr.Description = value
	return attr
}

// Tag declares an OpenAPI tag for a route operation.
func Tag(name string) *TagAttr {
	return &TagAttr{Name: name}
}

// TagAttr describes a route operation tag.
type TagAttr struct {
	Name string
}

// Summary declares the summary text for a route operation.
func Summary(summary string) *SummaryAttr {
	return &SummaryAttr{Summary: summary}
}

// SummaryAttr describes a route operation summary.
type SummaryAttr struct {
	Summary string
}

// Auth declares a security requirement for a route operation.
func Auth(name string, scopes ...string) *AuthAttr {
	return &AuthAttr{Name: name, Scopes: scopes}
}

// AuthAttr describes an OpenAPI security requirement.
type AuthAttr struct {
	Name   string
	Scopes []string
}

// Header creates an OpenAPI header parameter with a schema inferred from T.
func Header[T any](name string) *spec.Parameter {
	param := spec.HeaderParam(name)
	param.Schema = TypeToSchema(reflect.TypeFor[T]())
	return param
}

// Cookie creates an OpenAPI cookie parameter with a schema inferred from T.
func Cookie[T any](name string) *spec.Parameter {
	param := &spec.Parameter{ParamProps: spec.ParamProps{Name: name, In: "cookie"}}
	param.Schema = TypeToSchema(reflect.TypeFor[T]())
	return param
}
