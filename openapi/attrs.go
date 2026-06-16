package openapi

import (
	"reflect"

	"github.com/go-openapi/spec"
)

func Respond[T any](code int) *ResponseAttr {
	return &ResponseAttr{
		Code:        code,
		ContentType: "application/json",
		Type:        reflect.TypeFor[T](),
	}
}

func RespondFile(code int) *ResponseAttr {
	return &ResponseAttr{
		Code:        code,
		ContentType: "application/octet-stream",
		Schema: &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"file"},
			},
		},
	}
}

type ResponseAttr struct {
	Code        int
	ContentType string
	Type        reflect.Type
	Schema      *spec.Schema
	Description string
}

func (attr *ResponseAttr) WithContentType(value string) *ResponseAttr {
	attr.ContentType = value
	return attr
}

func (attr *ResponseAttr) WithDescription(value string) *ResponseAttr {
	attr.Description = value
	return attr
}

func Tag(name string) *TagAttr {
	return &TagAttr{Name: name}
}

type TagAttr struct {
	Name string
}

func Summary(summary string) *SummaryAttr {
	return &SummaryAttr{Summary: summary}
}

type SummaryAttr struct {
	Summary string
}
