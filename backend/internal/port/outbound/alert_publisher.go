package outbound

import "github.com/leveling/demonit/internal/domain"

// AlertPublisher broadcasts domain events to realtime subscribers (SSE/WebSocket).
type AlertPublisher interface {
	PublishDeviceOffline(event domain.DeviceOfflineEvent)
}
