package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/leveling/demonit/internal/adapter/outbound/realtime"
	"github.com/leveling/demonit/internal/config"
	"github.com/leveling/demonit/internal/port/inbound"
	"go.uber.org/zap"
)

// RouterDeps groups dependencies injected into the HTTP adapter layer.
type RouterDeps struct {
	DeviceService inbound.DeviceService
	Validate      *validator.Validate
	Logger        *zap.Logger
	Auth          config.AuthConfig
	CORS          config.CORSConfig
	Hub           *realtime.Hub
}

// NewRouter builds the Gin engine with all API routes registered.
//
// Auth model (intentionally not full user login):
//   - POST /heartbeat     → device API key  (machines / IoT agents)
//   - POST /devices       → admin API key   (dashboard register)
//   - GET  /devices*      → public read     (dashboard polling)
//   - GET  /events        → public SSE      (realtime offline alerts)
//   - GET  /healthz       → public
func NewRouter(deps RouterDeps) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware(deps.CORS.AllowedOrigins))
	r.Use(requestLogger(deps.Logger))

	heartbeat := NewHeartbeatHandler(deps.DeviceService, deps.Validate, deps.Logger)
	devices := NewDeviceHandler(deps.DeviceService, deps.Validate, deps.Logger)
	events := NewEventsHandler(deps.Hub, deps.Logger)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/heartbeat", requireAPIKey(HeaderDeviceAPIKey, deps.Auth.DeviceAPIKey), heartbeat.Handle)
		v1.POST("/devices", requireAPIKey(HeaderAdminAPIKey, deps.Auth.AdminAPIKey), devices.Create)
		v1.GET("/devices", devices.List)
		v1.GET("/devices/:id", devices.Get)
		v1.GET("/devices/:id/metrics", devices.ListMetrics)
		v1.GET("/events", events.Stream)
	}

	r.GET("/healthz", func(c *gin.Context) {
		JSONSuccess(c, 200, "ok", gin.H{"status": "healthy"})
	})

	return r
}

func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[o] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if _, ok := allowed[origin]; ok {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
				c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, "+HeaderDeviceAPIKey+", "+HeaderAdminAPIKey)
				c.Header("Access-Control-Max-Age", "86400")
			}
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func requestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		log.Info("http request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
