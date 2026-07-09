package openapi

import (
	"errors"

	"github.com/go-openapi/spec"
	"github.com/oesand/ino"
)

// GenerateOptions configures top-level OpenAPI document metadata.
type GenerateOptions struct {
	Info spec.InfoProps
	Tags []spec.Tag
}

// GenerateSchema builds an OpenAPI 2.0 schema from the routes registered in mux.
func GenerateSchema(mux *ino.Mux, options *GenerateOptions) (*spec.Swagger, error) {
	if mux == nil {
		return nil, errors.New("ino: mux is required")
	}

	var info spec.Info

	if options != nil {
		info.InfoProps = options.Info
	}

	if info.Title == "" {
		info.Title = "API"
	}

	if info.Version == "" {
		info.Version = "1.0.0"
	}

	swagger := newSwagger(&info)

	if options != nil {
		swagger.Tags = options.Tags
	}

	for route := range mux.Routes() {
		pattern, operation := newOperation(route)
		operation.FillConsumes(route, swagger)
		operation.ScanAttributes(route, swagger)

		var pathItem spec.PathItem
		operation.Apply(&pathItem, route.Method())
		swagger.Paths.Paths[pattern] = pathItem
	}

	return swagger.Base(), nil
}
