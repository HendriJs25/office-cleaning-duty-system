package middleware

import (
	"cleaning/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery // use if only in production raw query doesn't contain credential
		if raw != "" {
			path = path + "?" + raw
		}

		entry := logger.WithRequestID(GetRequestID(c)).WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       path,
			"status":     status,
			"latency_ms": latency.Milliseconds(),
			"client_ip":  c.ClientIP(),
		})

		entry.Info("request completed")
	}
}
