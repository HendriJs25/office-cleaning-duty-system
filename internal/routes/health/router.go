package health

import (
	healthhandler "cleaning/internal/handler/health"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, handler *healthhandler.Health) {
	router.GET("/health", handler.HealthCheck)
}
