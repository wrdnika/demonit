package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DeviceType enumerates supported hardware categories.
type DeviceType string

const (
	DeviceTypeATM    DeviceType = "ATM"
	DeviceTypeServer DeviceType = "SERVER"
	DeviceTypeLaptop DeviceType = "LAPTOP"
)

// DeviceStatus represents connectivity state.
type DeviceStatus string

const (
	DeviceStatusOnline  DeviceStatus = "ONLINE"
	DeviceStatusOffline DeviceStatus = "OFFLINE"
)

// Device is the aggregate root for a monitored IoT endpoint.
type Device struct {
	ID        uuid.UUID    `json:"id"`
	Name      string       `json:"name"`
	Type      DeviceType   `json:"type"`
	Status    DeviceStatus `json:"status"`
	LastSeen  time.Time    `json:"last_seen"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	// Latest metric sample (populated by list/detail queries when available).
	CPUUsage      *float64        `json:"cpu_usage,omitempty"`
	RAMUsage      *float64        `json:"ram_usage,omitempty"`
	StatusPayload json.RawMessage `json:"status_payload,omitempty"`
}

// RegisterDeviceInput is the command to enroll a new monitored device.
type RegisterDeviceInput struct {
	Name string     `json:"name" validate:"required,min=2,max=255"`
	Type DeviceType `json:"type" validate:"required"`
}

// UpdateDeviceInput is the command to rename/retarget an existing device.
type UpdateDeviceInput struct {
	Name string     `json:"name" validate:"required,min=2,max=255"`
	Type DeviceType `json:"type" validate:"required"`
}

// IsValidDeviceType reports whether t is a known DeviceType.
func IsValidDeviceType(t DeviceType) bool {
	switch t {
	case DeviceTypeATM, DeviceTypeServer, DeviceTypeLaptop:
		return true
	default:
		return false
	}
}
