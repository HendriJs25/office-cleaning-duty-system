package auth

import (
	userhandler "cleaning/internal/handler/user"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, handler *userhandler.Handler) {
	router.POST("/login", handler.Login)
}
