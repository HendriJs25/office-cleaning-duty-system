package auth

import (
	authHandler "cleaning/internal/handler/auth"
	rolehandler "cleaning/internal/handler/role"
	userhandler "cleaning/internal/handler/user"
	"cleaning/internal/middleware"
	roleroutes "cleaning/internal/routes/role"
	"cleaning/internal/routes/user"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup,
	authHandler *authHandler.Handler,
	roleHandler *rolehandler.Handler,
	userHandler *userhandler.Handler,
	authentication *middleware.Authentication,
	authorization *middleware.Authorization) {

	router.POST("/login", authHandler.Login)

	authenticated := router.Group("/auth")
	authenticated.Use(authentication.Handle())

	authenticated.POST("/logout", authHandler.Logout)

	roleroutes.Register(authenticated, roleHandler, authorization)
	user.Register(authenticated, userHandler, authorization)
}
