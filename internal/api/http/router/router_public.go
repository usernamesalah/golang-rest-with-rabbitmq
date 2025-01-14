package router

import (
	"tenant/internal/api/http/handler"
	"tenant/internal/container"

	"github.com/labstack/echo/v4"
)

func publicRouter(e *echo.Echo, hc *container.HandlerComponent) {

	// Tenant
	tenantHandler := handler.NewTenantHandler(hc)
	tenantRoute := e.Group("/tenants")
	{
		tenantRoute.POST("", tenantHandler.CreateTenant)
		tenantRoute.POST("/:clientID/process", tenantHandler.ProcessPayload)
		tenantRoute.DELETE("/:clientID", tenantHandler.DeleteTenant)
	}

}
