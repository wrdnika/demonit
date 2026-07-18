package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DeviceMetric is a single time-series sample from a device heartbeat.
type DeviceMetric struct {
	ID            uuid.UUID       `json:"id"`
	DeviceID      uuid.UUID       `json:"device_id"`
	CPUUsage      float64         `json:"cpu_usage"`
	RAMUsage      float64         `json:"ram_usage"`
	StatusPayload json.RawMessage `json:"status_payload"`
	Timestamp     time.Time       `json:"timestamp"`
}

// HeartbeatInput is the validated command payload for a device ping.
type HeartbeatInput struct {
	DeviceID      uuid.UUID       `json:"device_id" validate:"required"`
	CPUUsage      float64         `json:"cpu_usage" validate:"min=0,max=100"`
	RAMUsage      float64         `json:"ram_usage" validate:"min=0,max=100"`
	StatusPayload json.RawMessage `json:"status_payload" validate:"required"`
}
