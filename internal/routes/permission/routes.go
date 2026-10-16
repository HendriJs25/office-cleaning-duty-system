package permission

import (
	"cleaning/internal/constants"
	permissionhandler "cleaning/internal/handler/permission"
	"cleaning/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, handler *permissionhandler.Handler, authorization *middleware.Authorization) {
	router.GET("/permissions/options", authorization.RequirePermission(constants.PermissionRoleCreate), handler.GetAllPermissionsForOptions)
}
