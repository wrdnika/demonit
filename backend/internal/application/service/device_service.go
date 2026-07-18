package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/leveling/demonit/internal/domain"
	"github.com/leveling/demonit/internal/port/inbound"
	"github.com/leveling/demonit/internal/port/outbound"
	"go.uber.org/zap"
)

// Ensure compile-time compliance with the inbound port.
var _ inbound.DeviceService = (*DeviceService)(nil)

// DeviceService orchestrates device monitoring use cases.
type DeviceService struct {
	repo   outbound.DeviceRepository
	logger *zap.Logger
}

// NewDeviceService wires a DeviceService with its dependencies.
func NewDeviceService(repo outbound.DeviceRepository, logger *zap.Logger) *DeviceService {
	return &DeviceService{
		repo:   repo,
		logger: logger,
	}
}

// ProcessHeartbeat marks the device ONLINE, updates last_seen, and stores a metric sample.
func (s *DeviceService) ProcessHeartbeat(ctx context.Context, input domain.HeartbeatInput) error {
	if !json.Valid(input.StatusPayload) {
		return fmt.Errorf("%w: status_payload must be valid JSON", domain.ErrInvalidPayload)
	}

	now := time.Now().UTC()
	metric := &domain.DeviceMetric{
		ID:            uuid.New(),
		DeviceID:      input.DeviceID,
		CPUUsage:      input.CPUUsage,
		RAMUsage:      input.RAMUsage,
		StatusPayload: input.StatusPayload,
		Timestamp:     now,
	}

	if err := s.repo.RecordHeartbeat(ctx, input.DeviceID, metric); err != nil {
		return err
	}

	s.logger.Info("heartbeat recorded",
		zap.String("device_id", input.DeviceID.String()),
		zap.Float64("cpu_usage", input.CPUUsage),
		zap.Float64("ram_usage", input.RAMUsage),
	)
	return nil
}

// ListDevices returns all registered devices with current status.
func (s *DeviceService) ListDevices(ctx context.Context) ([]domain.Device, error) {
	return s.repo.FindAll(ctx)
}

// MarkStaleDevicesOffline flips ONLINE devices older than olderThanSeconds to OFFLINE
// and logs an alert for each transition. Returns the number of devices marked offline.
func (s *DeviceService) MarkStaleDevicesOffline(ctx context.Context, olderThanSeconds int) (int, error) {
	threshold := time.Now().UTC().Add(-time.Duration(olderThanSeconds) * time.Second)

	offline, err := s.repo.MarkOffline(ctx, threshold)
	if err != nil {
		return 0, err
	}

	for _, d := range offline {
		s.logger.Warn("ALERT: device went offline (deadman's switch)",
			zap.String("device_id", d.ID.String()),
			zap.String("device_name", d.Name),
			zap.String("device_type", string(d.Type)),
			zap.Time("last_seen", d.LastSeen),
			zap.Duration("stale_for", time.Since(d.LastSeen)),
		)
	}

	return len(offline), nil
}
