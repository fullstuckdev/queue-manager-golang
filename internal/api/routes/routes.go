package routes

import (
	"github.com/labstack/echo/v4"
	"queue-manager/internal/api/handlers"
)

func SetupRoutes(e *echo.Echo, tenantHandler *handlers.TenantHandler) {
	e.POST("/tenants", tenantHandler.CreateTenant)
	e.DELETE("/tenants/:clientID", tenantHandler.DeleteTenant)
	e.POST("/tenants/process", tenantHandler.ProcessPayload)
} 