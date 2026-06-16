package openapi

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-openapi/spec"
)

type User struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email,omitempty"`
	Active bool   `json:"active"`
}

func TestSchema(t *testing.T) {
	schema := structToSchema(User{})

	data, _ := json.MarshalIndent(schema, "", "  ")

	fmt.Println(string(data))
}

func main() {
	// Описание модели User
	userSchema := spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type: []string{"object"},
			Properties: map[string]spec.Schema{
				"id": {
					SchemaProps: spec.SchemaProps{
						Type:   []string{"integer"},
						Format: "int64",
					},
				},
				"name": {
					SchemaProps: spec.SchemaProps{
						Type: []string{"string"},
					},
				},
				"email": {
					SchemaProps: spec.SchemaProps{
						Type:   []string{"string"},
						Format: "email",
					},
				},
			},
			Required: []string{"id", "name"},
		},
	}

	// GET /users/{id}
	getUserOperation := &spec.Operation{
		OperationProps: spec.OperationProps{
			Summary: "Get user by id",
			ID:      "getUserById",
			Parameters: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name:     "id",
						In:       "path",
						Required: true,
					},
					SimpleSchema: spec.SimpleSchema{
						Type:   "integer",
						Format: "int64",
					},
				},
			},
			Responses: &spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					StatusCodeResponses: map[int]spec.Response{
						200: {
							ResponseProps: spec.ResponseProps{
								Description: "OK",
								Schema:      spec.RefSchema("#/definitions/User"),
							},
						},
					},
				},
			},
		},
	}

	// Swagger документ
	swagger := spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger: "2.0",
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Title:       "Example API",
					Version:     "1.0.0",
					Description: "Simple API example",
				},
			},
			Paths: &spec.Paths{
				Paths: map[string]spec.PathItem{
					"/users/{id}": {
						PathItemProps: spec.PathItemProps{
							Get: getUserOperation,
						},
					},
				},
			},
			Definitions: spec.Definitions{
				"User": userSchema,
			},
		},
	}

	data, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
