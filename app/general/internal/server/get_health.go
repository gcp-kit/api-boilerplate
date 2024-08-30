package server

import (
	"api-boilerplate/app/general/internal/interfaces/openapi"
	"context"

	"github.com/cockroachdb/errors"
)

func (s Server) GetHealth(ctx context.Context, request openapi.GetHealthRequestObject) (openapi.GetHealthResponseObject, error) {
	err := s.props.HealthCheckUsecase.HealthCheck(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call HealthCheckUsecase.HealthCheck")
	}
	res := openapi.GetHealth200JSONResponse{
		openapi.HealthCheckJSONResponse{
			Status: "ok",
		},
	}
	return res, nil
}
