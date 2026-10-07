package middleware

import (
	"cleaning/internal/common/response"
	"io"
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func HandlePanic() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, recovered any) {
		slog.Error(
			"panic recovered",
			"error", recovered,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"stack", string(debug.Stack()),
		)

		response.InternalServerError(c)
		c.Abort()
	})
}
