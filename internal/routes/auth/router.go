package auth

import (
	authHandler "cleaning/internal/handler/auth"
	employeehandler "cleaning/internal/handler/employee"
	officehandler "cleaning/internal/handler/office"
	rolehandler "cleaning/internal/handler/role"
	userhandler "cleaning/internal/handler/user"
	verificationtokenhandler "cleaning/internal/handler/verificationtoken"
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
	verificationTokenHandler *verificationtokenhandler.Handler,
	authentication *middleware.Authentication,
	authorization *middleware.Authorization) {

	auth := router.Group("/auth")

	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authentication.Handle(), authHandler.Logout)

	authenticated := router.Group("")
	authenticated.Use(authentication.Handle())

	roleroutes.Register(authenticated, roleHandler, authorization)
	userroutes.Register(authenticated, userHandler, authorization)
	officeroutes.Register(authenticated, officeHandler, authorization)
	employee.Register(authenticated, employeeHandler, verificationTokenHandler, authorization)
}
