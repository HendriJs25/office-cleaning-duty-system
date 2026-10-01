package role

import (
	"cleaning/internal/constants"
	rolehandler "cleaning/internal/handler/role"
	"cleaning/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, handler *rolehandler.Handler, authorization *middleware.Authorization) {
	router.GET("/roles", authorization.RequirePermission(constants.PermissionRoleRead), handler.GetAll)
	router.GET("/roles/options", authorization.RequirePermission(constants.PermissionRoleRead), handler.GetAllActiveRoles)
}
