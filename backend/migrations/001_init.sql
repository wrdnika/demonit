-- devices: monitored IoT endpoints
CREATE TABLE IF NOT EXISTS devices (
    id         UUID PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    type       VARCHAR(32)  NOT NULL CHECK (type IN ('ATM', 'SERVER', 'LAPTOP')),
    status     VARCHAR(16)  NOT NULL DEFAULT 'OFFLINE' CHECK (status IN ('ONLINE', 'OFFLINE')),
    last_seen  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_devices_status ON devices (status);
CREATE INDEX IF NOT EXISTS idx_devices_last_seen ON devices (last_seen);
CREATE INDEX IF NOT EXISTS idx_devices_type ON devices (type);

-- device_metrics: time-series heartbeat samples
CREATE TABLE IF NOT EXISTS device_metrics (
    id             UUID PRIMARY KEY,
    device_id      UUID           NOT NULL REFERENCES devices (id) ON DELETE CASCADE,
    cpu_usage      DOUBLE PRECISION NOT NULL,
    ram_usage      DOUBLE PRECISION NOT NULL,
    status_payload JSONB          NOT NULL DEFAULT '{}'::jsonb,
    timestamp      TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_device_metrics_device_id ON device_metrics (device_id);
CREATE INDEX IF NOT EXISTS idx_device_metrics_timestamp ON device_metrics (timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_device_metrics_device_ts ON device_metrics (device_id, timestamp DESC);
