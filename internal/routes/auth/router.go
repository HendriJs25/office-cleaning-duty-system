package auth

import (
	authHandler "cleaning/internal/handler/auth"
	rolehandler "cleaning/internal/handler/role"
	"cleaning/internal/middleware"
	roleroutes "cleaning/internal/routes/role"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, authHandler *authHandler.Handler, roleHandler *rolehandler.Handler, authentication *middleware.Authentication) {
	router.POST("/login", authHandler.Login)

	authenticated := router.Group("/auth")
	authenticated.Use(authentication.Handle())
	authenticated.POST("/logout", authHandler.Logout)

	roleroutes.Register(authenticated, roleHandler)
}
