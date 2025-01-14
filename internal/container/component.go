package container

import (
	"tenant/infrastructure/config"
	"tenant/internal/repository"
	"tenant/internal/service/messaging"
	"tenant/internal/usecase"
)

type HandlerComponent struct {
	Config *config.Config

	// Usecase
	TenantUsecase usecase.TenantUsecase
}

func NewHandlerComponent(sc *SharedComponent) *HandlerComponent {

	// Library
	mq := messaging.NewRabbitMQ(sc.RabbitMQConn)

	tenantRepo := repository.NewTenantRepository(sc.DB)
	tenantUsecase := usecase.NewTenantUsecase(tenantRepo, mq)

	return &HandlerComponent{
		Config: sc.Conf,

		// Usecase
		TenantUsecase: tenantUsecase,
	}
}
