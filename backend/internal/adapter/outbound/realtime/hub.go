package realtime

import (
	"encoding/json"
	"sync"

	"github.com/leveling/demonit/internal/domain"
	"github.com/leveling/demonit/internal/port/outbound"
)

// Ensure compile-time compliance.
var _ outbound.AlertPublisher = (*Hub)(nil)

// Hub is an in-process fan-out bus for SSE subscribers (single-node).
type Hub struct {
	mu      sync.RWMutex
	clients map[chan []byte]struct{}
}

// NewHub constructs an empty realtime hub.
func NewHub() *Hub {
	return &Hub{clients: make(map[chan []byte]struct{})}
}

// Subscribe registers a buffered client channel. Caller must Unsubscribe.
func (h *Hub) Subscribe() chan []byte {
	ch := make(chan []byte, 8)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

// Unsubscribe removes and closes the client channel.
func (h *Hub) Unsubscribe(ch chan []byte) {
	h.mu.Lock()
	if _, ok := h.clients[ch]; ok {
		delete(h.clients, ch)
		close(ch)
	}
	h.mu.Unlock()
}

// PublishDeviceOffline fans out a JSON SSE payload to all subscribers.
func (h *Hub) PublishDeviceOffline(event domain.DeviceOfflineEvent) {
	payload, err := json.Marshal(map[string]any{
		"type": "device_offline",
		"data": event,
	})
	if err != nil {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.clients {
		select {
		case ch <- payload:
		default:
			// Slow consumer: drop rather than block the publisher.
		}
	}
}

// ClientCount returns the number of active SSE subscribers (for tests/ops).
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
