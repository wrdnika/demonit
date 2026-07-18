package outbound

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/leveling/demonit/internal/domain"
)

// DeviceRepository is the outbound port for device persistence.
// Implementations live in adapter/outbound (e.g. PostgreSQL).
type DeviceRepository interface {
	FindAll(ctx context.Context) ([]domain.Device, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Device, error)
	UpdateHeartbeat(ctx context.Context, deviceID uuid.UUID, lastSeen time.Time) error
	MarkOffline(ctx context.Context, threshold time.Time) ([]domain.Device, error)
	InsertMetric(ctx context.Context, metric *domain.DeviceMetric) error
	// RecordHeartbeat atomically updates device status + inserts a metric row.
	RecordHeartbeat(ctx context.Context, deviceID uuid.UUID, metric *domain.DeviceMetric) error
}
