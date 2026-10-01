package office

import (
	"cleaning/internal/constants"
	officehandler "cleaning/internal/handler/office"
	"cleaning/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, handler *officehandler.Handler, authorization *middleware.Authorization) {
	router.POST("/offices", authorization.RequirePermission(constants.PermissionOfficeCreate), handler.Create)
	router.GET("/offices/options", authorization.RequirePermission(constants.PermissionOfficeRead), handler.GetAllActiveOffices)
}
