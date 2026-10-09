package employee

import (
	"cleaning/internal/constants"
	employeehandler "cleaning/internal/handler/employee"
	"cleaning/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, handler *employeehandler.Handler, authorization *middleware.Authorization) {
	router.POST("/employees", authorization.RequirePermission(constants.PermissionEmployeeCreate), handler.Create)
	router.GET("/employees/options", authorization.RequirePermission(constants.PermissionEmployeeRead), handler.GetAllAssignableEmployees)
}
