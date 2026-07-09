package openapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-openapi/spec"
	"github.com/oesand/mo"
	httpSwagger "github.com/swaggo/http-swagger"
)

// RegisterOptions configures the OpenAPI JSON and Swagger UI route patterns.
type RegisterOptions struct {
	OpenApiPattern string
	SwaggerPattern string
}

// Register mounts routes that serve the generated OpenAPI JSON and Swagger UI.
func Register(mux *mo.Mux, schema *spec.Swagger, options *RegisterOptions) error {
	var openApiPattern, swaggerPattern string

	if options != nil {
		openApiPattern = options.OpenApiPattern
		swaggerPattern = options.SwaggerPattern
	}

	if openApiPattern == "" {
		openApiPattern = "/openapi.json"
	}

	if swaggerPattern == "" {
		swaggerPattern = "/swagger/{*}"
	}

	if schema == nil {
		return errors.New("mo: schema is required")
	}

	schemaData, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	mux.Include(
		mo.Get(openApiPattern, mo.F(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(schemaData)
		})),
		mo.Get(swaggerPattern, httpSwagger.Handler(
			httpSwagger.URL(openApiPattern),
		)),
	)
	return nil
}
