package openapi

import (
	"reflect"
	"strings"

	"github.com/go-openapi/spec"
)

// TypeToSchema converts a Go reflection type into an OpenAPI schema.
func TypeToSchema(typ reflect.Type) *spec.Schema {
	switch typ.Kind() {
	case reflect.String:
		return spec.StringProperty()
	case reflect.Bool:
		return spec.BoolProperty()
	case reflect.Int8:
		return spec.Int8Property()
	case reflect.Int16:
		return spec.Int16Property()
	case reflect.Int, reflect.Int32:
		return spec.Int32Property()
	case reflect.Int64:
		return spec.Int64Property()
	case reflect.Float32:
		return spec.Float32Property()
	case reflect.Float64:
		return spec.Float64Property()
	case reflect.Slice, reflect.Array:
		return spec.ArrayProperty(TypeToSchema(typ.Elem()))
	case reflect.Map:
		schema := new(spec.Schema).Typed("object", "")

		key := typ.Key()
		if key.Kind() != reflect.String {
			// Non-string keys: not standard in OpenAPI 2.0.
			return schema
		}

		elem := typ.Elem()
		schema.AdditionalProperties = &spec.SchemaOrBool{Schema: TypeToSchema(elem)}
		return schema
	case reflect.Pointer:
		return TypeToSchema(typ.Elem())
	case reflect.Struct:
		return structToSchema(typ)
	default:
		return &spec.Schema{}
	}
}

// structToSchema converts exported and unexported struct fields into object properties using json tags.
func structToSchema(structType reflect.Type) *spec.Schema {
	properties := map[string]spec.Schema{}
	var required []string

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		name := field.Name
		omitempty := false

		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				name = parts[0]
			}

			for _, p := range parts[1:] {
				if p == "omitempty" {
					omitempty = true
				}
			}
		}

		schema := TypeToSchema(field.Type)

		properties[name] = *schema

		if !omitempty {
			required = append(required, name)
		}
	}

	return &spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type:       []string{"object"},
			Properties: properties,
			Required:   required,
		},
	}
}
