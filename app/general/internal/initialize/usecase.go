package initialize

import "api-boilerplate/app/internal/usecases/healthcheck"

// Usecases - usecases
type Usecases struct {
	HealthCheckUsecase *healthcheck.Usecase
}

func NewUsecases() Usecases {
	healthCheckUsecase := healthcheck.NewUsecase()
	return Usecases{
		HealthCheckUsecase: healthCheckUsecase,
	}
}
