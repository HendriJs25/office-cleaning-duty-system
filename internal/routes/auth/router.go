package auth

import (
	authHandler "cleaning/internal/handler/auth"
	employeehandler "cleaning/internal/handler/employee"
	officehandler "cleaning/internal/handler/office"
	permissionhandler "cleaning/internal/handler/permission"
	rolehandler "cleaning/internal/handler/role"
	userhandler "cleaning/internal/handler/user"
	verificationtokenhandler "cleaning/internal/handler/verificationtoken"
	"cleaning/internal/middleware"
	"cleaning/internal/routes/employee"
	officeroutes "cleaning/internal/routes/office"
	"cleaning/internal/routes/permission"
	roleroutes "cleaning/internal/routes/role"
	userroutes "cleaning/internal/routes/user"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup,
	authHandler *authHandler.Handler,
	permissionHandler *permissionhandler.Handler,
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
	auth.POST("/registration-tokens/verify", verificationTokenHandler.VerifyRegistrationToken)
	auth.POST("/register", authHandler.Register)

	authenticated := router.Group("")
	authenticated.Use(authentication.Handle())

	permission.Register(authenticated, permissionHandler, authorization)
	roleroutes.Register(authenticated, roleHandler, authorization)
	userroutes.Register(authenticated, userHandler, authorization)
	officeroutes.Register(authenticated, officeHandler, authorization)
	employee.Register(authenticated, employeeHandler, verificationTokenHandler, authorization)
}
