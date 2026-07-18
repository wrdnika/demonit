package postgres

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/leveling/demonit/internal/domain"
	"gorm.io/datatypes"
)

// DeviceModel is the GORM persistence model for devices.
type DeviceModel struct {
	ID        uuid.UUID          `gorm:"type:uuid;primaryKey"`
	Name      string             `gorm:"type:varchar(255);not null"`
	Type      domain.DeviceType  `gorm:"type:varchar(32);not null;index"`
	Status    domain.DeviceStatus `gorm:"type:varchar(16);not null;index;default:'OFFLINE'"`
	LastSeen  time.Time          `gorm:"not null;index"`
	CreatedAt time.Time          `gorm:"not null"`
	UpdatedAt time.Time          `gorm:"not null"`
}

func (DeviceModel) TableName() string { return "devices" }

func (m DeviceModel) ToDomain() domain.Device {
	return domain.Device{
		ID:        m.ID,
		Name:      m.Name,
		Type:      m.Type,
		Status:    m.Status,
		LastSeen:  m.LastSeen,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// DeviceMetricModel is the GORM persistence model for time-series metrics.
type DeviceMetricModel struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	DeviceID      uuid.UUID      `gorm:"type:uuid;not null;index"`
	CPUUsage      float64        `gorm:"type:double precision;not null"`
	RAMUsage      float64        `gorm:"type:double precision;not null"`
	StatusPayload datatypes.JSON `gorm:"type:jsonb;not null"`
	Timestamp     time.Time      `gorm:"not null;index"`
}

func (DeviceMetricModel) TableName() string { return "device_metrics" }

func (m DeviceMetricModel) ToDomain() domain.DeviceMetric {
	return domain.DeviceMetric{
		ID:            m.ID,
		DeviceID:      m.DeviceID,
		CPUUsage:      m.CPUUsage,
		RAMUsage:      m.RAMUsage,
		StatusPayload: json.RawMessage(m.StatusPayload),
		Timestamp:     m.Timestamp,
	}
}

func metricFromDomain(m *domain.DeviceMetric) DeviceMetricModel {
	payload := m.StatusPayload
	if payload == nil {
		payload = json.RawMessage(`{}`)
	}
	return DeviceMetricModel{
		ID:            m.ID,
		DeviceID:      m.DeviceID,
		CPUUsage:      m.CPUUsage,
		RAMUsage:      m.RAMUsage,
		StatusPayload: datatypes.JSON(payload),
		Timestamp:     m.Timestamp,
	}
}
