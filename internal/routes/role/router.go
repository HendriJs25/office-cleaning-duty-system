package role

import (
	"cleaning/internal/constants"
	rolehandler "cleaning/internal/handler/role"
	"cleaning/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, handler *rolehandler.Handler, authorization *middleware.Authorization) {
	router.GET("/roles", authorization.RequirePermission(constants.PermissionRoleRead), handler.GetAll)
	router.POST("/roles", authorization.RequirePermission(constants.PermissionRoleCreate), authorization.RequirePermission(constants.PermissionRoleAssignPermission), handler.Create)
	router.GET("/roles/:id", authorization.RequirePermission(constants.PermissionRoleRead), handler.GetRoleDetail)
	router.GET("/roles/options", authorization.RequirePermission(constants.PermissionRoleRead), handler.GetAllActiveRoles)
	router.PATCH("/roles/:id/deactivate", authorization.RequirePermission(constants.PermissionRoleDeActivate), handler.DeactivateRole)
}
