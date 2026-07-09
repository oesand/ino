package openapi

import (
	"reflect"
	"testing"
)

type schemaTestStruct struct {
	ID       int64             `json:"id"`
	Name     string            `json:"name"`
	Email    string            `json:"email,omitempty"`
	Ignored  string            `json:"-"`
	Metadata map[string]int32  `json:"metadata"`
	Flags    []bool            `json:"flags"`
	Pointer  *schemaNestedType `json:"pointer"`
}

type schemaNestedType struct {
	Value string `json:"value"`
}

func TestTypeToSchemaBasicTypes(t *testing.T) {
	tests := []struct {
		name   string
		typ    reflect.Type
		kind   string
		format string
	}{
		{"string", reflect.TypeFor[string](), "string", ""},
		{"bool", reflect.TypeFor[bool](), "boolean", ""},
		{"int8", reflect.TypeFor[int8](), "integer", "int8"},
		{"int16", reflect.TypeFor[int16](), "integer", "int16"},
		{"int", reflect.TypeFor[int](), "integer", "int32"},
		{"int32", reflect.TypeFor[int32](), "integer", "int32"},
		{"int64", reflect.TypeFor[int64](), "integer", "int64"},
		{"float32", reflect.TypeFor[float32](), "number", "float"},
		{"float64", reflect.TypeFor[float64](), "number", "double"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := TypeToSchema(tt.typ)
			if len(schema.Type) != 1 || schema.Type[0] != tt.kind {
				t.Fatalf("expected type %s, got %v", tt.kind, schema.Type)
			}
			if schema.Format != tt.format {
				t.Fatalf("expected format %s, got %s", tt.format, schema.Format)
			}
		})
	}
}

func TestTypeToSchemaCollectionTypes(t *testing.T) {
	arraySchema := TypeToSchema(reflect.TypeFor[[]int64]())
	if arraySchema.Type[0] != "array" {
		t.Fatalf("expected array schema, got %v", arraySchema.Type)
	}
	if arraySchema.Items.Schema.Type[0] != "integer" || arraySchema.Items.Schema.Format != "int64" {
		t.Fatalf("expected int64 array item schema, got type=%v format=%s", arraySchema.Items.Schema.Type, arraySchema.Items.Schema.Format)
	}

	mapSchema := TypeToSchema(reflect.TypeFor[map[string]bool]())
	if mapSchema.Type[0] != "object" {
		t.Fatalf("expected object schema, got %v", mapSchema.Type)
	}
	if mapSchema.AdditionalProperties == nil || mapSchema.AdditionalProperties.Schema.Type[0] != "boolean" {
		t.Fatalf("expected boolean additional properties, got %#v", mapSchema.AdditionalProperties)
	}

	unsupportedKeyMap := TypeToSchema(reflect.TypeFor[map[int]string]())
	if unsupportedKeyMap.AdditionalProperties != nil {
		t.Fatalf("expected no additional properties for non-string map keys, got %#v", unsupportedKeyMap.AdditionalProperties)
	}
}

func TestTypeToSchemaPointerAndUnsupportedTypes(t *testing.T) {
	pointerSchema := TypeToSchema(reflect.TypeFor[*string]())
	if pointerSchema.Type[0] != "string" {
		t.Fatalf("expected pointer to resolve to string schema, got %v", pointerSchema.Type)
	}

	unsupportedSchema := TypeToSchema(reflect.TypeFor[chan int]())
	if len(unsupportedSchema.Type) != 0 {
		t.Fatalf("expected empty schema for unsupported type, got %v", unsupportedSchema.Type)
	}
}

func TestStructToSchemaUsesJSONTagsAndRequiredFields(t *testing.T) {
	schema := structToSchema(reflect.TypeFor[schemaTestStruct]())

	if schema.Type[0] != "object" {
		t.Fatalf("expected object schema, got %v", schema.Type)
	}
	if _, ok := schema.Properties["Ignored"]; ok {
		t.Fatal("expected json:- field to be skipped")
	}
	if _, ok := schema.Properties["ignored"]; ok {
		t.Fatal("expected json:- field to be skipped")
	}
	if schema.Properties["id"].Format != "int64" {
		t.Fatalf("expected id int64 format, got %s", schema.Properties["id"].Format)
	}
	if schema.Properties["metadata"].AdditionalProperties.Schema.Format != "int32" {
		t.Fatalf("expected metadata int32 values, got %#v", schema.Properties["metadata"].AdditionalProperties.Schema)
	}
	if schema.Properties["flags"].Items.Schema.Type[0] != "boolean" {
		t.Fatalf("expected flags boolean items, got %#v", schema.Properties["flags"].Items.Schema.Type)
	}
	if schema.Properties["pointer"].Type[0] != "object" {
		t.Fatalf("expected pointer struct to resolve to object, got %v", schema.Properties["pointer"].Type)
	}

	required := map[string]bool{}
	for _, name := range schema.Required {
		required[name] = true
	}
	for _, name := range []string{"id", "name", "metadata", "flags", "pointer"} {
		if !required[name] {
			t.Fatalf("expected %s to be required in %#v", name, schema.Required)
		}
	}
	if required["email"] {
		t.Fatalf("expected omitempty field email not to be required: %#v", schema.Required)
	}
}
