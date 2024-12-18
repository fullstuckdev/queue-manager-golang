package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	
	"queue-manager/internal/models"
	"queue-manager/internal/queue"
	"queue-manager/internal/repository"
)

type TenantHandler struct {
	tenantRepo *repository.TenantRepository
	rmq        *queue.RabbitMQ
	logger     *logrus.Logger
}

func NewTenantHandler(repo *repository.TenantRepository, rmq *queue.RabbitMQ, logger *logrus.Logger) *TenantHandler {
	return &TenantHandler{
		tenantRepo: repo,
		rmq:        rmq,
		logger:     logger,
	}
}

func (h *TenantHandler) CreateTenant(c echo.Context) error {
	var req models.CreateTenantRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	tenant, err := h.tenantRepo.Create(c.Request().Context(), &req)
	if err != nil {
		h.logger.Errorf("Failed to create tenant: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create tenant"})
	}

	if err := h.rmq.CreateQueue(req.ClientID); err != nil {
		h.logger.Errorf("Failed to create queue: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create queue"})
	}

	if err := h.rmq.StartConsumer(req.ClientID); err != nil {
		h.logger.Errorf("Failed to start consumer: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start consumer"})
	}

	return c.JSON(http.StatusCreated, tenant)
}

func (h *TenantHandler) DeleteTenant(c echo.Context) error {
	clientID := c.Param("clientID")

	tenant, err := h.tenantRepo.GetByClientID(c.Request().Context(), clientID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found"})
	}

	h.rmq.StopConsumer(clientID)

	if err := h.tenantRepo.Delete(c.Request().Context(), clientID); err != nil {
		h.logger.Errorf("Failed to delete tenant: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete tenant"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *TenantHandler) ProcessPayload(c echo.Context) error {
	var req models.ProcessPayloadRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	_, err := h.tenantRepo.GetByClientID(c.Request().Context(), req.ClientID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found"})
	}

	payloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload format"})
	}

	if err := h.rmq.PublishMessage(req.ClientID, payloadBytes); err != nil {
		h.logger.Errorf("Failed to publish message: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to publish message"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "Message published successfully"})
} 