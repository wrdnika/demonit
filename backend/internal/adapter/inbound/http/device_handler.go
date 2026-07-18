package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leveling/demonit/internal/port/inbound"
	"go.uber.org/zap"
)

// DeviceHandler serves device listing endpoints.
type DeviceHandler struct {
	service inbound.DeviceService
	logger  *zap.Logger
}

// NewDeviceHandler constructs a DeviceHandler with explicit DI.
func NewDeviceHandler(service inbound.DeviceService, logger *zap.Logger) *DeviceHandler {
	return &DeviceHandler{
		service: service,
		logger:  logger,
	}
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
