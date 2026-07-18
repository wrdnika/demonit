package http

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/leveling/demonit/internal/adapter/outbound/realtime"
	"go.uber.org/zap"
)

// EventsHandler streams server-sent events for realtime alerts.
type EventsHandler struct {
	hub    *realtime.Hub
	logger *zap.Logger
}

// NewEventsHandler constructs an SSE handler.
func NewEventsHandler(hub *realtime.Hub, logger *zap.Logger) *EventsHandler {
	return &EventsHandler{hub: hub, logger: logger}
}

// Stream handles GET /api/v1/events (text/event-stream).
func (h *EventsHandler) Stream(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	ch := h.hub.Subscribe()
	defer h.hub.Unsubscribe(ch)

	// Initial comment keeps intermediaries from buffering an empty stream.
	_, _ = fmt.Fprintf(c.Writer, ": connected\n\n")
	c.Writer.Flush()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			return false
		case payload, ok := <-ch:
			if !ok {
				return false
			}
			if _, err := fmt.Fprintf(w, "event: device_offline\ndata: %s\n\n", payload); err != nil {
				h.logger.Debug("sse write failed", zap.Error(err))
				return false
			}
			return true
		}
	})
}
