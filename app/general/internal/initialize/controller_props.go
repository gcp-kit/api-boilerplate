package initialize

import (
	"api-boilerplate/app/general/internal/config"
	"api-boilerplate/app/general/internal/interfaces/props"
)

func NewControllerProps(cfg *config.Config, usecases Usecases) *props.ControllerProps {
	cp := &props.ControllerProps{
		HealthCheckUsecase: usecases.HealthCheckUsecase,
	}
	return cp
}
