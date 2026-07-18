package domain

import (
	"time"

	"github.com/google/uuid"
)

// DeviceOfflineEvent is emitted when the deadman's switch flips a device OFFLINE.
type DeviceOfflineEvent struct {
	DeviceID   uuid.UUID `json:"device_id"`
	DeviceName string    `json:"device_name"`
	DeviceType string    `json:"device_type"`
	LastSeen   time.Time `json:"last_seen"`
	OccurredAt time.Time `json:"occurred_at"`
}
