package worker

import (
	"context"
	"sync"
	"time"

	"github.com/leveling/demonit/internal/port/inbound"
	"go.uber.org/zap"
)

// DeadmanSwitch periodically marks devices offline when heartbeats stop arriving.
// It is safe to run as a single long-lived goroutine per process.
type DeadmanSwitch struct {
	service        inbound.DeviceService
	interval       time.Duration
	offlineTimeout time.Duration
	logger         *zap.Logger

	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewDeadmanSwitch constructs a worker. Call Start to begin ticking.
func NewDeadmanSwitch(
	service inbound.DeviceService,
	interval time.Duration,
	offlineTimeout time.Duration,
	logger *zap.Logger,
) *DeadmanSwitch {
	return &DeadmanSwitch{
		service:        service,
		interval:       interval,
		offlineTimeout: offlineTimeout,
		logger:         logger,
	}
}

// Start launches the background ticker loop. It is non-blocking.
// Calling Start more than once without Stop is undefined; wire once from main.
func (w *DeadmanSwitch) Start(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)
	w.cancel = cancel

	w.wg.Add(1)
	go w.loop(ctx)
	w.logger.Info("deadman switch started",
		zap.Duration("interval", w.interval),
		zap.Duration("offline_timeout", w.offlineTimeout),
	)
}

// Stop cancels the loop and waits for the goroutine to exit.
func (w *DeadmanSwitch) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
	w.wg.Wait()
	w.logger.Info("deadman switch stopped")
}

func (w *DeadmanSwitch) loop(ctx context.Context) {
	defer w.wg.Done()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Run once immediately so cold starts don't wait a full interval.
	w.tick(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}

func (w *DeadmanSwitch) tick(ctx context.Context) {
	// Bound each tick so a slow DB query cannot pile up work.
	tickCtx, cancel := context.WithTimeout(ctx, w.interval)
	defer cancel()

	seconds := int(w.offlineTimeout.Seconds())
	if seconds < 1 {
		seconds = 1
	}

	count, err := w.service.MarkStaleDevicesOffline(tickCtx, seconds)
	if err != nil {
		w.logger.Error("deadman switch tick failed", zap.Error(err))
		return
	}

	if count > 0 {
		w.logger.Info("deadman switch marked devices offline", zap.Int("count", count))
	}
}
