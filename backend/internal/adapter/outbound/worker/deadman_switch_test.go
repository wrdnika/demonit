package worker_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/leveling/demonit/internal/adapter/outbound/worker"
	"github.com/leveling/demonit/internal/domain"
	"go.uber.org/zap"
)

type stubService struct {
	calls []int
}

func (s *stubService) RegisterDevice(ctx context.Context, input domain.RegisterDeviceInput) (*domain.Device, error) {
	return nil, nil
}
func (s *stubService) UpdateDevice(ctx context.Context, id uuid.UUID, input domain.UpdateDeviceInput) (*domain.Device, error) {
	return nil, nil
}
func (s *stubService) DeleteDevice(ctx context.Context, id uuid.UUID) error { return nil }
func (s *stubService) ProcessHeartbeat(ctx context.Context, input domain.HeartbeatInput) error {
	return nil
}
func (s *stubService) ListDevices(ctx context.Context) ([]domain.Device, error) { return nil, nil }
func (s *stubService) GetDevice(ctx context.Context, id uuid.UUID) (*domain.Device, error) {
	return nil, domain.ErrDeviceNotFound
}
func (s *stubService) ListDeviceMetrics(ctx context.Context, id uuid.UUID, limit int) ([]domain.DeviceMetric, error) {
	return nil, nil
}
func (s *stubService) MarkStaleDevicesOffline(ctx context.Context, olderThanSeconds int) (int, error) {
	s.calls = append(s.calls, olderThanSeconds)
	return 1, nil
}

func TestDeadmanSwitch_TicksAndStops(t *testing.T) {
	svc := &stubService{}
	w := worker.NewDeadmanSwitch(svc, 40*time.Millisecond, 30*time.Second, zap.NewNop())

	ctx, cancel := context.WithCancel(context.Background())
	w.Start(ctx)

	deadline := time.Now().Add(300 * time.Millisecond)
	for time.Now().Before(deadline) {
		if len(svc.calls) >= 2 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	cancel()
	w.Stop()

	if len(svc.calls) < 2 {
		t.Fatalf("expected at least 2 ticks (immediate + ticker), got %d", len(svc.calls))
	}
	for _, seconds := range svc.calls {
		if seconds != 30 {
			t.Fatalf("expected offline timeout seconds=30, got %d", seconds)
		}
	}
}
