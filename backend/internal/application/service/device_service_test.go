package service_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/leveling/demonit/internal/application/service"
	"github.com/leveling/demonit/internal/domain"
	"go.uber.org/zap"
)

type mockRepo struct {
	heartbeatCalls int
	lastMetric     *domain.DeviceMetric
	offline        []domain.Device
	markErr        error
	heartbeatErr   error
}

func (m *mockRepo) Create(ctx context.Context, device *domain.Device) error { return nil }
func (m *mockRepo) Update(ctx context.Context, device *domain.Device) error { return nil }
func (m *mockRepo) Delete(ctx context.Context, id uuid.UUID) error          { return nil }
func (m *mockRepo) FindAll(ctx context.Context) ([]domain.Device, error)     { return nil, nil }
func (m *mockRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Device, error) {
	return nil, domain.ErrDeviceNotFound
}
func (m *mockRepo) MarkOffline(ctx context.Context, threshold time.Time) ([]domain.Device, error) {
	if m.markErr != nil {
		return nil, m.markErr
	}
	return m.offline, nil
}
func (m *mockRepo) RecordHeartbeat(ctx context.Context, deviceID uuid.UUID, metric *domain.DeviceMetric) error {
	m.heartbeatCalls++
	m.lastMetric = metric
	return m.heartbeatErr
}
func (m *mockRepo) ListMetrics(ctx context.Context, deviceID uuid.UUID, limit int) ([]domain.DeviceMetric, error) {
	return nil, nil
}

type mockAlerts struct {
	events []domain.DeviceOfflineEvent
}

func (m *mockAlerts) PublishDeviceOffline(event domain.DeviceOfflineEvent) {
	m.events = append(m.events, event)
}

func TestProcessHeartbeat_RecordsMetric(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewDeviceService(repo, nil, zap.NewNop())
	id := uuid.New()

	err := svc.ProcessHeartbeat(context.Background(), domain.HeartbeatInput{
		DeviceID:      id,
		CPUUsage:      12.5,
		RAMUsage:      40,
		StatusPayload: json.RawMessage(`{"temperature":36.5}`),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.heartbeatCalls != 1 {
		t.Fatalf("expected 1 RecordHeartbeat call, got %d", repo.heartbeatCalls)
	}
	if repo.lastMetric == nil || repo.lastMetric.DeviceID != id {
		t.Fatalf("metric not recorded for device")
	}
}

func TestProcessHeartbeat_RejectsInvalidJSON(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewDeviceService(repo, nil, zap.NewNop())

	err := svc.ProcessHeartbeat(context.Background(), domain.HeartbeatInput{
		DeviceID:      uuid.New(),
		CPUUsage:      1,
		RAMUsage:      1,
		StatusPayload: json.RawMessage(`{bad`),
	})
	if err == nil {
		t.Fatal("expected error for invalid JSON payload")
	}
	if repo.heartbeatCalls != 0 {
		t.Fatal("RecordHeartbeat must not be called on invalid payload")
	}
}

func TestMarkStaleDevicesOffline_PublishesAlerts(t *testing.T) {
	id := uuid.New()
	repo := &mockRepo{
		offline: []domain.Device{{
			ID:       id,
			Name:     "ATM-1",
			Type:     domain.DeviceTypeATM,
			Status:   domain.DeviceStatusOffline,
			LastSeen: time.Now().UTC().Add(-time.Minute),
		}},
	}
	alerts := &mockAlerts{}
	svc := service.NewDeviceService(repo, alerts, zap.NewNop())

	n, err := svc.MarkStaleDevicesOffline(context.Background(), 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 offline device, got %d", n)
	}
	if len(alerts.events) != 1 {
		t.Fatalf("expected 1 published alert, got %d", len(alerts.events))
	}
	if alerts.events[0].DeviceID != id || alerts.events[0].DeviceName != "ATM-1" {
		t.Fatalf("unexpected alert payload: %+v", alerts.events[0])
	}
}
