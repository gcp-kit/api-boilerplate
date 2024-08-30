package server

import (
	"api-boilerplate/app/general/internal/interfaces/openapi"
	"api-boilerplate/app/general/internal/interfaces/props"
)

// Server - server
type Server struct {
	props *props.ControllerProps
}

// NewServer - constructor
func NewServer(cp *props.ControllerProps) openapi.StrictServerInterface {
	return &Server{
		props: cp,
	}
}
