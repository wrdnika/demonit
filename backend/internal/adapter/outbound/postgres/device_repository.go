package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/leveling/demonit/internal/domain"
	"github.com/leveling/demonit/internal/port/outbound"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Ensure compile-time compliance with the outbound port.
var _ outbound.DeviceRepository = (*DeviceRepository)(nil)

// DeviceRepository is the PostgreSQL adapter for DeviceRepository.
type DeviceRepository struct {
	db *gorm.DB
}

// NewDeviceRepository constructs a repository backed by db.
func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) FindAll(ctx context.Context) ([]domain.Device, error) {
	var rows []DeviceModel
	if err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("find all devices: %w", err)
	}

	devices := make([]domain.Device, 0, len(rows))
	for _, row := range rows {
		devices = append(devices, row.ToDomain())
	}
	return devices, nil
}

func (r *DeviceRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Device, error) {
	var row DeviceModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrDeviceNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find device by id: %w", err)
	}
	d := row.ToDomain()
	return &d, nil
}

func (r *DeviceRepository) UpdateHeartbeat(ctx context.Context, deviceID uuid.UUID, lastSeen time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&DeviceModel{}).
		Where("id = ?", deviceID).
		Updates(map[string]any{
			"status":     domain.DeviceStatusOnline,
			"last_seen":  lastSeen,
			"updated_at": lastSeen,
		})
	if result.Error != nil {
		return fmt.Errorf("update heartbeat: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrDeviceNotFound
	}
	return nil
}

func (r *DeviceRepository) InsertMetric(ctx context.Context, metric *domain.DeviceMetric) error {
	row := metricFromDomain(metric)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("insert metric: %w", err)
	}
	return nil
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
		// Lock matching ONLINE rows to avoid duplicate alerts under concurrent workers.
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
