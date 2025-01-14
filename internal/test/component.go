package test

import (
	"tenant/infrastructure/config"
	"tenant/internal/test/mockrepository"
	"tenant/internal/test/mockusecase"
	"testing"
)

type MockComponent struct {
	Config           *config.Config
	TenantRepository *mockrepository.TenantRepository
	TenantUsecase    *mockusecase.TenantUsecase
}

func InitMockComponent(t *testing.T) *MockComponent {
	return &MockComponent{
		Config:           &config.Config{},
		TenantRepository: mockrepository.NewTenantRepository(t),
		TenantUsecase:    mockusecase.NewTenantUsecase(t),
	}
}
