package middleware

import (
	"cleaning/internal/common/response"
	"cleaning/internal/constants"
	"io"
	"log/slog"
	"net/http"
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

		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status:  constants.Error,
			Message: "サーバー内部でエラーが発生しました。",
		})

	})
}
