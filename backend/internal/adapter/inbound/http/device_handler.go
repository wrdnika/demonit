package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/leveling/demonit/internal/domain"
	"github.com/leveling/demonit/internal/port/inbound"
	"go.uber.org/zap"
)

// DeviceHandler serves device CRUD and metrics endpoints.
type DeviceHandler struct {
	service  inbound.DeviceService
	validate *validator.Validate
	logger   *zap.Logger
}

// NewDeviceHandler constructs a DeviceHandler with explicit DI.
func NewDeviceHandler(service inbound.DeviceService, validate *validator.Validate, logger *zap.Logger) *DeviceHandler {
	return &DeviceHandler{
		service:  service,
		validate: validate,
		logger:   logger,
	}
}

type registerDeviceRequest struct {
	Name string            `json:"name" validate:"required,min=2,max=255"`
	Type domain.DeviceType `json:"type" validate:"required"`
}

// Create handles POST /api/v1/devices.
func (h *DeviceHandler) Create(c *gin.Context) {
	var req registerDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONError(c, http.StatusBadRequest, "INVALID_JSON", "request body must be valid JSON", nil)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		JSONError(c, http.StatusBadRequest, "VALIDATION_ERROR", "request validation failed", ValidationDetails(err))
		return
	}

	device, err := h.service.RegisterDevice(c.Request.Context(), domain.RegisterDeviceInput{
		Name: req.Name,
		Type: req.Type,
	})
	if err != nil {
		h.logger.Error("register device failed", zap.Error(err))
		switch {
		case errors.Is(err, domain.ErrInvalidDeviceType):
			JSONError(c, http.StatusBadRequest, "INVALID_DEVICE_TYPE", "type must be ATM, SERVER, or LAPTOP", nil)
		default:
			JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to register device", nil)
		}
		return
	}

	JSONSuccess(c, http.StatusCreated, "device registered", device)
}

// List handles GET /api/v1/devices.
func (h *DeviceHandler) List(c *gin.Context) {
	devices, err := h.service.ListDevices(c.Request.Context())
	if err != nil {
		h.logger.Error("list devices failed", zap.Error(err))
		JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list devices", nil)
		return
	}

	JSONSuccess(c, http.StatusOK, "devices retrieved", devices)
}

// Get handles GET /api/v1/devices/:id.
func (h *DeviceHandler) Get(c *gin.Context) {
	id, ok := parseDeviceID(c)
	if !ok {
		return
	}

	device, err := h.service.GetDevice(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrDeviceNotFound) {
			JSONError(c, http.StatusNotFound, "DEVICE_NOT_FOUND", "device not found", nil)
			return
		}
		h.logger.Error("get device failed", zap.Error(err), zap.String("device_id", id.String()))
		JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get device", nil)
		return
	}

	JSONSuccess(c, http.StatusOK, "device retrieved", device)
}

// ListMetrics handles GET /api/v1/devices/:id/metrics.
func (h *DeviceHandler) ListMetrics(c *gin.Context) {
	id, ok := parseDeviceID(c)
	if !ok {
		return
	}

	limit := 50
	if raw := c.Query("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			limit = n
		}
	}

	metrics, err := h.service.ListDeviceMetrics(c.Request.Context(), id, limit)
	if err != nil {
		if errors.Is(err, domain.ErrDeviceNotFound) {
			JSONError(c, http.StatusNotFound, "DEVICE_NOT_FOUND", "device not found", nil)
			return
		}
		h.logger.Error("list metrics failed", zap.Error(err), zap.String("device_id", id.String()))
		JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list metrics", nil)
		return
	}

	JSONSuccess(c, http.StatusOK, "metrics retrieved", metrics)
}

func parseDeviceID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		JSONError(c, http.StatusBadRequest, "VALIDATION_ERROR", "id must be a valid UUID", map[string]string{
			"id": "must be a valid UUID",
		})
		return uuid.Nil, false
	}
	return id, true
}
