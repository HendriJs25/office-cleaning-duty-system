package auth

import (
	authHandler "cleaning/internal/handler/auth"
	employeehandler "cleaning/internal/handler/employee"
	officehandler "cleaning/internal/handler/office"
	rolehandler "cleaning/internal/handler/role"
	userhandler "cleaning/internal/handler/user"
	"cleaning/internal/middleware"
	"cleaning/internal/routes/employee"
	officeroutes "cleaning/internal/routes/office"
	roleroutes "cleaning/internal/routes/role"
	userroutes "cleaning/internal/routes/user"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup,
	authHandler *authHandler.Handler,
	roleHandler *rolehandler.Handler,
	userHandler *userhandler.Handler,
	officeHandler *officehandler.Handler,
	employeeHandler *employeehandler.Handler,
	authentication *middleware.Authentication,
	authorization *middleware.Authorization) {

	router.POST("/login", authHandler.Login)

	authenticated := router.Group("/auth")
	authenticated.Use(authentication.Handle())

	authenticated.POST("/logout", authHandler.Logout)

	roleroutes.Register(authenticated, roleHandler, authorization)
	userroutes.Register(authenticated, userHandler, authorization)
	officeroutes.Register(authenticated, officeHandler, authorization)
	employee.Register(authenticated, employeeHandler, authorization)
}
