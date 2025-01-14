package handler

import (
	"tenant/internal/container"
	"tenant/pkg/api"
	"tenant/pkg/derrors"

	"tenant/internal/api/http/handler/request"
	"tenant/internal/usecase"

	"github.com/labstack/echo/v4"
)

type (
	tenantHandler struct {
		tenantUsecase usecase.TenantUsecase
	}

	TenantHandler interface {
		CreateTenant(c echo.Context) error
		DeleteTenant(c echo.Context) error
		ProcessPayload(c echo.Context) error
	}
)

func NewTenantHandler(hc *container.HandlerComponent) TenantHandler {
	return &tenantHandler{tenantUsecase: hc.TenantUsecase}
}

// CreateTenant handles tenant creation requests
// Create Tenant
// @Summary Create Tenant
// @Description Create Tenant
// @Tags tenant
// @ID create-tenant
// @Produce json
// @Param user body request.CreateTenantRequest true "create tenant payload"
// @Success 200 {object} map[string]string
// @Router /tenant [post]
func (h *tenantHandler) CreateTenant(c echo.Context) error {
	var req request.CreateTenantRequest
	if err := c.Bind(&req); err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	if err := c.Validate(req); err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	ctx := c.Request().Context()
	tenant, err := h.tenantUsecase.CreateTenant(ctx, req.Name)
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	return api.ResponseOK(c, map[string]interface{}{
		"message": "Tenant created successfully",
		"tenant":  tenant,
	})
}

// DeleteTenant handles tenant deletion requests
// Delete Tenant
// @Summary Delete Tenant
// @Description Delete Tenant
// @Tags tenant
// @ID delete-tenant
// @Produce json
// @Param clientID path string true "clientID"
// @Success 200 {object} map[string]string
// @Router /tenant/{clientID} [delete]
func (h *tenantHandler) DeleteTenant(c echo.Context) error {
	clientID := c.Param("clientID")
	if clientID == "" {
		return api.RenderErrorResponse(c, c.Request(), derrors.New(derrors.InvalidArgument, "client_id is required"))
	}

	ctx := c.Request().Context()
	err := h.tenantUsecase.DeleteTenant(ctx, clientID)
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	return api.ResponseOK(c, map[string]interface{}{
		"message": "Tenant deleted successfully",
	})
}

// ProcessPayload handles requests to publish payload to a tenant's RabbitMQ queue
// Process Tenant
// @Summary Process Tenant
// @Description Process Tenant
// @Tags tenant
// @ID process-tenant
// @Produce json
// @Param clientID path string true "clientID"
// @Param user body request.ProcessPayloadRequest true "process tenant payload"
// @Success 200 {object} map[string]string
// @Router /tenant/{clientID}/process [post]
func (h *tenantHandler) ProcessPayload(c echo.Context) error {
	var req request.ProcessPayloadRequest
	clientID := c.Param("clientID")
	if clientID == "" {
		return api.RenderErrorResponse(c, c.Request(), derrors.New(derrors.InvalidArgument, "client_id is required"))
	}

	if err := c.Bind(&req); err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	if err := c.Validate(req); err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	ctx := c.Request().Context()
	err := h.tenantUsecase.ProcessPayload(ctx, clientID, req.Payload)
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	return api.ResponseOK(c, map[string]interface{}{
		"message": "Payload processed successfully",
	})
}
