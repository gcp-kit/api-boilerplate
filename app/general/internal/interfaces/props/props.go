package props

import "api-boilerplate/app/internal/usecases/healthcheck"

// ControllerProps - controller props
type ControllerProps struct {
	HealthCheckUsecase *healthcheck.Usecase
}
