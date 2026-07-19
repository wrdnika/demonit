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
	repo     outbound.DeviceRepository
	alerts   outbound.AlertPublisher
	logger   *zap.Logger
}

// NewDeviceService wires a DeviceService with its dependencies.
// alerts may be nil; offline events are then only logged.
func NewDeviceService(repo outbound.DeviceRepository, alerts outbound.AlertPublisher, logger *zap.Logger) *DeviceService {
	return &DeviceService{
		repo:   repo,
		alerts: alerts,
		logger: logger,
	}
}

// RegisterDevice enrolls a new device in OFFLINE state until the first heartbeat.
func (s *DeviceService) RegisterDevice(ctx context.Context, input domain.RegisterDeviceInput) (*domain.Device, error) {
	if !domain.IsValidDeviceType(input.Type) {
		return nil, domain.ErrInvalidDeviceType
	}

	now := time.Now().UTC()
	device := &domain.Device{
		ID:        uuid.New(),
		Name:      input.Name,
		Type:      input.Type,
		Status:    domain.DeviceStatusOffline,
		LastSeen:  now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, device); err != nil {
		return nil, err
	}

	s.logger.Info("device registered",
		zap.String("device_id", device.ID.String()),
		zap.String("name", device.Name),
		zap.String("type", string(device.Type)),
	)
	return device, nil
}

// UpdateDevice changes name/type of an existing device.
func (s *DeviceService) UpdateDevice(ctx context.Context, id uuid.UUID, input domain.UpdateDeviceInput) (*domain.Device, error) {
	if !domain.IsValidDeviceType(input.Type) {
		return nil, domain.ErrInvalidDeviceType
	}

	device, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	device.Name = input.Name
	device.Type = input.Type
	device.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, device); err != nil {
		return nil, err
	}

	s.logger.Info("device updated",
		zap.String("device_id", device.ID.String()),
		zap.String("name", device.Name),
		zap.String("type", string(device.Type)),
	)
	return device, nil
}

// DeleteDevice removes a device and cascades metrics via FK.
func (s *DeviceService) DeleteDevice(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.logger.Info("device deleted", zap.String("device_id", id.String()))
	return nil
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

// ListDevices returns all registered devices with latest metric samples when available.
func (s *DeviceService) ListDevices(ctx context.Context) ([]domain.Device, error) {
	return s.repo.FindAll(ctx)
}

// GetDevice returns a single device with its latest metrics.
func (s *DeviceService) GetDevice(ctx context.Context, id uuid.UUID) (*domain.Device, error) {
	return s.repo.FindByID(ctx, id)
}

// ListDeviceMetrics returns recent time-series samples for a device.
func (s *DeviceService) ListDeviceMetrics(ctx context.Context, id uuid.UUID, limit int) ([]domain.DeviceMetric, error) {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	return s.repo.ListMetrics(ctx, id, limit)
}

// MarkStaleDevicesOffline flips ONLINE devices older than olderThanSeconds to OFFLINE,
// logs alerts, and publishes realtime events for SSE subscribers.
func (s *DeviceService) MarkStaleDevicesOffline(ctx context.Context, olderThanSeconds int) (int, error) {
	threshold := time.Now().UTC().Add(-time.Duration(olderThanSeconds) * time.Second)

	offline, err := s.repo.MarkOffline(ctx, threshold)
	if err != nil {
		return 0, err
	}

	now := time.Now().UTC()
	for _, d := range offline {
		s.logger.Warn("ALERT: device went offline (deadman's switch)",
			zap.String("device_id", d.ID.String()),
			zap.String("device_name", d.Name),
			zap.String("device_type", string(d.Type)),
			zap.Time("last_seen", d.LastSeen),
			zap.Duration("stale_for", time.Since(d.LastSeen)),
		)

		if s.alerts != nil {
			s.alerts.PublishDeviceOffline(domain.DeviceOfflineEvent{
				DeviceID:   d.ID,
				DeviceName: d.Name,
				DeviceType: string(d.Type),
				LastSeen:   d.LastSeen,
				OccurredAt: now,
			})
		}
	}

	return len(offline), nil
}
