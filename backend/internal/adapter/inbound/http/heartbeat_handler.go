package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/leveling/demonit/internal/domain"
	"github.com/leveling/demonit/internal/port/inbound"
	"go.uber.org/zap"
)

// HeartbeatRequest is the HTTP DTO for POST /api/v1/heartbeat.
type HeartbeatRequest struct {
	DeviceID      string          `json:"device_id" validate:"required,uuid"`
	CPUUsage      float64         `json:"cpu_usage" validate:"min=0,max=100"`
	RAMUsage      float64         `json:"ram_usage" validate:"min=0,max=100"`
	StatusPayload json.RawMessage `json:"status_payload" validate:"required"`
}

// HeartbeatHandler handles device heartbeat ingestion.
type HeartbeatHandler struct {
	service  inbound.DeviceService
	validate *validator.Validate
	logger   *zap.Logger
}

// NewHeartbeatHandler constructs a HeartbeatHandler with explicit DI.
func NewHeartbeatHandler(service inbound.DeviceService, validate *validator.Validate, logger *zap.Logger) *HeartbeatHandler {
	return &HeartbeatHandler{
		service:  service,
		validate: validate,
		logger:   logger,
	}
}

// Handle processes POST /api/v1/heartbeat.
// Updates device last_seen + status=ONLINE and inserts a device_metrics row.
func (h *HeartbeatHandler) Handle(c *gin.Context) {
	var req HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONError(c, http.StatusBadRequest, "INVALID_JSON", "request body must be valid JSON", nil)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		JSONError(c, http.StatusBadRequest, "VALIDATION_ERROR", "request validation failed", ValidationDetails(err))
		return
	}

	deviceID, err := uuid.Parse(req.DeviceID)
	if err != nil {
		JSONError(c, http.StatusBadRequest, "VALIDATION_ERROR", "device_id must be a valid UUID", map[string]string{
			"DeviceID": "must be a valid UUID",
		})
		return
	}

	if !json.Valid(req.StatusPayload) {
		JSONError(c, http.StatusBadRequest, "VALIDATION_ERROR", "status_payload must be valid JSON", map[string]string{
			"StatusPayload": "must be valid JSON",
		})
		return
	}

	input := domain.HeartbeatInput{
		DeviceID:      deviceID,
		CPUUsage:      req.CPUUsage,
		RAMUsage:      req.RAMUsage,
		StatusPayload: req.StatusPayload,
	}

	if err := h.service.ProcessHeartbeat(c.Request.Context(), input); err != nil {
		h.logger.Error("heartbeat processing failed",
			zap.String("device_id", req.DeviceID),
			zap.Error(err),
		)

		switch {
		case errors.Is(err, domain.ErrDeviceNotFound):
			JSONError(c, http.StatusNotFound, "DEVICE_NOT_FOUND", "device not found", nil)
		case errors.Is(err, domain.ErrInvalidPayload):
			JSONError(c, http.StatusBadRequest, "INVALID_PAYLOAD", err.Error(), nil)
		default:
			JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to process heartbeat", nil)
		}
		return
	}

	JSONSuccess(c, http.StatusOK, "heartbeat accepted", gin.H{
		"device_id": deviceID,
		"status":    domain.DeviceStatusOnline,
	})
}
