package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Worker   WorkerConfig
	Auth     AuthConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Host string
	Port string
	Env  string
}

type DatabaseConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type WorkerConfig struct {
	Interval       time.Duration
	OfflineTimeout time.Duration
}

// AuthConfig holds shared secrets for machine and admin write access.
// Device endpoints (heartbeat) use DeviceAPIKey.
// Dashboard mutations (register device) use AdminAPIKey.
// Full user/password login is intentionally deferred — IoT agents cannot log in.
type AuthConfig struct {
	DeviceAPIKey string
	AdminAPIKey  string
}

type CORSConfig struct {
	AllowedOrigins []string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			DSN:             getEnv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=demonit port=5432 sslmode=disable TimeZone=UTC"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 2*time.Minute),
		},
		Worker: WorkerConfig{
			Interval:       getEnvAsDuration("WORKER_INTERVAL", 30*time.Second),
			OfflineTimeout: getEnvAsDuration("OFFLINE_TIMEOUT", 30*time.Second),
		},
		Auth: AuthConfig{
			DeviceAPIKey: getEnv("DEVICE_API_KEY", ""),
			AdminAPIKey:  getEnv("ADMIN_API_KEY", ""),
		},
		CORS: CORSConfig{
			AllowedOrigins: splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		},
	}

	if cfg.Database.DSN == "" {
		return nil, fmt.Errorf("DATABASE_DSN is required")
	}
	if cfg.Auth.DeviceAPIKey == "" {
		return nil, fmt.Errorf("DEVICE_API_KEY is required")
	}
	if cfg.Auth.AdminAPIKey == "" {
		return nil, fmt.Errorf("ADMIN_API_KEY is required")
	}

	return cfg, nil
}

func (c ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
