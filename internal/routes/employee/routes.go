package employee

import (
	"cleaning/internal/constants"
	employeehandler "cleaning/internal/handler/employee"
	verificationtokenhandler "cleaning/internal/handler/verificationtoken"
	"cleaning/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, handler *employeehandler.Handler, verificationTokenHandler *verificationtokenhandler.Handler, authorization *middleware.Authorization) {
	router.GET("/employees", authorization.RequirePermission(constants.PermissionEmployeeRead), handler.GetAll)
	router.POST("/employees", authorization.RequirePermission(constants.PermissionEmployeeCreate), handler.Create)
	router.GET("/employees/options", authorization.RequirePermission(constants.PermissionEmployeeRead), handler.GetAllAssignableEmployees)
	router.PATCH("/employees/:uuid", authorization.RequirePermission(constants.PermissionEmployeeUpdate), handler.UpdateEmployeeByAdmin)
	router.POST("/employees/:uuid/registration-tokens", authorization.RequirePermission(constants.PermissionEmployeeInvitationCreate), verificationTokenHandler.CreateRegistrationToken)
}
