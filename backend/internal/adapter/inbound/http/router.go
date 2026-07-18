package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/leveling/demonit/internal/port/inbound"
	"go.uber.org/zap"
)

// RouterDeps groups dependencies injected into the HTTP adapter layer.
type RouterDeps struct {
	DeviceService inbound.DeviceService
	Validate      *validator.Validate
	Logger        *zap.Logger
}

// NewRouter builds the Gin engine with all API routes registered.
func NewRouter(deps RouterDeps) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())
	r.Use(requestLogger(deps.Logger))

	heartbeat := NewHeartbeatHandler(deps.DeviceService, deps.Validate, deps.Logger)
	devices := NewDeviceHandler(deps.DeviceService, deps.Logger)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/heartbeat", heartbeat.Handle)
		v1.GET("/devices", devices.List)
	}

	r.GET("/healthz", func(c *gin.Context) {
		JSONSuccess(c, 200, "ok", gin.H{"status": "healthy"})
	})

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*"
		}

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

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
