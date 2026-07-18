package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	inboundhttp "github.com/leveling/demonit/internal/adapter/inbound/http"
	"github.com/leveling/demonit/internal/adapter/outbound/postgres"
	"github.com/leveling/demonit/internal/adapter/outbound/realtime"
	"github.com/leveling/demonit/internal/adapter/outbound/worker"
	"github.com/leveling/demonit/internal/application/service"
	"github.com/leveling/demonit/internal/config"
	"github.com/leveling/demonit/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		panic("config: " + err.Error())
	}

	log, err := logger.New(cfg.Server.Env)
	if err != nil {
		panic("logger: " + err.Error())
	}
	defer func() { _ = log.Sync() }()

	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		log.Fatal("database connection failed", zap.Error(err))
	}
	defer func() {
		if err := postgres.Close(db); err != nil {
			log.Error("database close failed", zap.Error(err))
		}
	}()

	// AutoMigrate only in local/dev. Production must apply migrations/001_init.sql.
	if isDevEnv(cfg.Server.Env) {
		if err := postgres.AutoMigrate(db); err != nil {
			log.Fatal("database migration failed", zap.Error(err))
		}
		log.Info("AutoMigrate applied (development)")
	} else {
		log.Info("skipping AutoMigrate; ensure SQL migrations are applied",
			zap.String("env", cfg.Server.Env),
		)
	}

	// --- Dependency injection (composition root) ---
	hub := realtime.NewHub()
	deviceRepo := postgres.NewDeviceRepository(db)
	deviceService := service.NewDeviceService(deviceRepo, hub, log)
	validate := validator.New()

	deadman := worker.NewDeadmanSwitch(
		deviceService,
		cfg.Worker.Interval,
		cfg.Worker.OfflineTimeout,
		log,
	)

	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	deadman.Start(rootCtx)

	router := inboundhttp.NewRouter(inboundhttp.RouterDeps{
		DeviceService: deviceService,
		Validate:      validate,
		Logger:        log,
		Auth:          cfg.Auth,
		CORS:          cfg.CORS,
		Hub:           hub,
	})

	srv := &http.Server{
		Addr:              cfg.Server.Addr(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		// WriteTimeout must be 0 for long-lived SSE connections.
		WriteTimeout: 0,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info("HTTP server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Info("shutdown signal received", zap.String("signal", sig.String()))

	rootCancel()
	deadman.Stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP shutdown error", zap.Error(err))
	}

	log.Info("server stopped cleanly")
}

func isDevEnv(env string) bool {
	switch strings.ToLower(strings.TrimSpace(env)) {
	case "", "development", "dev", "local":
		return true
	default:
		return false
	}
}
