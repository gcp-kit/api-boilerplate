package healthcheck

import "context"

// Usecase - usecase
type Usecase struct{}

// NewUsecase - usecase constructor
func NewUsecase() *Usecase {
	return &Usecase{}
}

// HealthCheck - health check
func (u *Usecase) HealthCheck(ctx context.Context) error {
	return nil
}
