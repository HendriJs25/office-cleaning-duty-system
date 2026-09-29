package user

import (
	"cleaning/internal/constants"
	userhandler "cleaning/internal/handler/user"
	"cleaning/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, handler *userhandler.Handler, authorization *middleware.Authorization) {
	router.GET("/users", authorization.RequirePermission(constants.PermissionUserRead), handler.GetAll)
	router.POST("/users", authorization.RequirePermission(constants.PermissionUserCreate), handler.Create)
}
