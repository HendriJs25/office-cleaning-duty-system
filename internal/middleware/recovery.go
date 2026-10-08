package middleware

import (
	"cleaning/internal/common/response"
	"cleaning/internal/logger"
	"io"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, recovered any) {
		logger.WithContext(c.Request.Context()).WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"error":  recovered,
			"stack":  string(debug.Stack()),
		}).Error("panic recovered")

		response.InternalServerError(c)
		c.Abort()
	})
}
