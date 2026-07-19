package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/leveling/demonit/internal/domain"
	"github.com/leveling/demonit/internal/port/outbound"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Ensure compile-time compliance with the outbound port.
var _ outbound.DeviceRepository = (*DeviceRepository)(nil)

// DeviceRepository is the PostgreSQL adapter for DeviceRepository.
type DeviceRepository struct {
	db *gorm.DB
}

// deviceWithMetrics is a scan target for devices joined with latest metrics.
type deviceWithMetrics struct {
	DeviceModel
	CPUUsage      *float64       `gorm:"column:cpu_usage"`
	RAMUsage      *float64       `gorm:"column:ram_usage"`
	StatusPayload datatypes.JSON `gorm:"column:status_payload"`
}

// NewDeviceRepository constructs a repository backed by db.
func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Create(ctx context.Context, device *domain.Device) error {
	row := DeviceModel{
		ID:        device.ID,
		Name:      device.Name,
		Type:      device.Type,
		Status:    device.Status,
		LastSeen:  device.LastSeen,
		CreatedAt: device.CreatedAt,
		UpdatedAt: device.UpdatedAt,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("create device: %w", err)
	}
	return nil
}

func (r *DeviceRepository) Update(ctx context.Context, device *domain.Device) error {
	result := r.db.WithContext(ctx).
		Model(&DeviceModel{}).
		Where("id = ?", device.ID).
		Updates(map[string]any{
			"name":       device.Name,
			"type":       device.Type,
			"updated_at": device.UpdatedAt,
		})
	if result.Error != nil {
		return fmt.Errorf("update device: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrDeviceNotFound
	}
	return nil
}

func (r *DeviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&DeviceModel{})
	if result.Error != nil {
		return fmt.Errorf("delete device: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrDeviceNotFound
	}
	return nil
}

func (d deviceWithMetrics) toDomainDevice() domain.Device {
	out := d.ToDomain()
	out.CPUUsage = d.CPUUsage
	out.RAMUsage = d.RAMUsage
	if len(d.StatusPayload) > 0 && string(d.StatusPayload) != "null" {
		out.StatusPayload = json.RawMessage(d.StatusPayload)
	}
	return out
}

func (r *DeviceRepository) FindAll(ctx context.Context) ([]domain.Device, error) {
	var rows []deviceWithMetrics

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			d.id, d.name, d.type, d.status, d.last_seen, d.created_at, d.updated_at,
			m.cpu_usage, m.ram_usage, m.status_payload
		FROM devices d
		LEFT JOIN LATERAL (
			SELECT cpu_usage, ram_usage, status_payload
			FROM device_metrics
			WHERE device_id = d.id
			ORDER BY timestamp DESC
			LIMIT 1
		) m ON TRUE
		ORDER BY d.name ASC
	`).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("find all devices: %w", err)
	}

	devices := make([]domain.Device, 0, len(rows))
	for _, row := range rows {
		devices = append(devices, row.toDomainDevice())
	}
	return devices, nil
}

func (r *DeviceRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Device, error) {
	var row deviceWithMetrics
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			d.id, d.name, d.type, d.status, d.last_seen, d.created_at, d.updated_at,
			m.cpu_usage, m.ram_usage, m.status_payload
		FROM devices d
		LEFT JOIN LATERAL (
			SELECT cpu_usage, ram_usage, status_payload
			FROM device_metrics
			WHERE device_id = d.id
			ORDER BY timestamp DESC
			LIMIT 1
		) m ON TRUE
		WHERE d.id = ?
		LIMIT 1
	`, id).Scan(&row).Error
	if err != nil {
		return nil, fmt.Errorf("find device by id: %w", err)
	}
	if row.ID == uuid.Nil {
		return nil, domain.ErrDeviceNotFound
	}

	d := row.toDomainDevice()
	return &d, nil
}

func (r *DeviceRepository) ListMetrics(ctx context.Context, deviceID uuid.UUID, limit int) ([]domain.DeviceMetric, error) {
	var rows []DeviceMetricModel
	if err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("timestamp DESC").
		Limit(limit).
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list metrics: %w", err)
	}

	out := make([]domain.DeviceMetric, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.ToDomain())
	}
	return out, nil
}

// RecordHeartbeat updates the device to ONLINE and inserts a metric in one transaction.
func (r *DeviceRepository) RecordHeartbeat(ctx context.Context, deviceID uuid.UUID, metric *domain.DeviceMetric) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := metric.Timestamp

		result := tx.Model(&DeviceModel{}).
			Where("id = ?", deviceID).
			Updates(map[string]any{
				"status":     domain.DeviceStatusOnline,
				"last_seen":  now,
				"updated_at": now,
			})
		if result.Error != nil {
			return fmt.Errorf("update device heartbeat: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return domain.ErrDeviceNotFound
		}

		row := metricFromDomain(metric)
		if err := tx.Create(&row).Error; err != nil {
			return fmt.Errorf("insert device metric: %w", err)
		}
		return nil
	})
}

// MarkOffline sets status=OFFLINE for devices whose last_seen is before threshold.
// Returns the devices that transitioned so callers can emit alerts.
func (r *DeviceRepository) MarkOffline(ctx context.Context, threshold time.Time) ([]domain.Device, error) {
	var stale []DeviceModel

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("status = ? AND last_seen < ?", domain.DeviceStatusOnline, threshold).
			Find(&stale).Error; err != nil {
			return fmt.Errorf("select stale devices: %w", err)
		}
		if len(stale) == 0 {
			return nil
		}

		ids := make([]uuid.UUID, 0, len(stale))
		for _, d := range stale {
			ids = append(ids, d.ID)
		}

		now := time.Now().UTC()
		if err := tx.Model(&DeviceModel{}).
			Where("id IN ?", ids).
			Updates(map[string]any{
				"status":     domain.DeviceStatusOffline,
				"updated_at": now,
			}).Error; err != nil {
			return fmt.Errorf("mark devices offline: %w", err)
		}

		for i := range stale {
			stale[i].Status = domain.DeviceStatusOffline
			stale[i].UpdatedAt = now
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	out := make([]domain.Device, 0, len(stale))
	for _, row := range stale {
		out = append(out, row.ToDomain())
	}
	return out, nil
}
