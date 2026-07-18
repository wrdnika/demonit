package inbound

import (
	"context"

	"github.com/leveling/demonit/internal/domain"
)

// DeviceService is the inbound port (use-case API) consumed by HTTP adapters
// and background workers.
type DeviceService interface {
	ProcessHeartbeat(ctx context.Context, input domain.HeartbeatInput) error
	ListDevices(ctx context.Context) ([]domain.Device, error)
	MarkStaleDevicesOffline(ctx context.Context, olderThanSeconds int) (int, error)
}
