package domain

import (
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
