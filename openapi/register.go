package openapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-openapi/spec"
	"github.com/oesand/ino"
	httpSwagger "github.com/swaggo/http-swagger"
)

type RegisterOptions struct {
	OpenApiPattern string
	SwaggerPattern string
}

func Register(mux *ino.Mux, schema *spec.Swagger, options RegisterOptions) error {
	if options.OpenApiPattern == "" {
		options.OpenApiPattern = "/openapi.json"
	}

	if options.SwaggerPattern == "" {
		options.SwaggerPattern = "/swagger/{*}"
	}

	if schema == nil {
		return errors.New("ino: schema is required")
	}

	schemaData, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	mux.Include(
		ino.Get(options.OpenApiPattern, ino.F(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(schemaData)
		})),
		ino.Get(options.SwaggerPattern, httpSwagger.Handler(
			httpSwagger.URL(options.OpenApiPattern),
		)),
	)
	return nil
}

// change name to "mo"
