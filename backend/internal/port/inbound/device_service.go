package inbound

import (
	"context"

	"github.com/google/uuid"
	"github.com/leveling/demonit/internal/domain"
)

// DeviceService is the inbound port (use-case API) consumed by HTTP adapters
// and background workers.
type DeviceService interface {
	RegisterDevice(ctx context.Context, input domain.RegisterDeviceInput) (*domain.Device, error)
	ProcessHeartbeat(ctx context.Context, input domain.HeartbeatInput) error
	ListDevices(ctx context.Context) ([]domain.Device, error)
	GetDevice(ctx context.Context, id uuid.UUID) (*domain.Device, error)
	ListDeviceMetrics(ctx context.Context, id uuid.UUID, limit int) ([]domain.DeviceMetric, error)
	MarkStaleDevicesOffline(ctx context.Context, olderThanSeconds int) (int, error)
}
