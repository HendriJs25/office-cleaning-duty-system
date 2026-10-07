package bootstrap

import (
	"cleaning/internal/config"

	"github.com/gin-gonic/gin"
)

func SetupAppMode(cfg config.App) {
	if cfg.IsDebugMode {
		gin.SetMode(gin.DebugMode)
		return
	}
	gin.SetMode(gin.ReleaseMode)
}
