package initialize

import (
	"github.com/cockroachdb/errors"
	"github.com/getkin/kin-openapi/openapi3"
)

// ExternalDependencies - external dependencies
type ExternalDependencies struct {
	openapi *openapi3.T
}

// NewExternalDependencies - initialize external dependencies
func NewExternalDependencies() (*ExternalDependencies, error) {
	ed := new(ExternalDependencies)

	{
		loader := openapi3.NewLoader()
		spec, err := loader.LoadFromData(OpenAPISpecBin)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get openapi specs")
		}
		ed.openapi = spec
	}

	return ed, nil
}

// OpenAPISpec - get openapi specs
func (e *ExternalDependencies) OpenAPISpec() *openapi3.T {
	return e.openapi
}
