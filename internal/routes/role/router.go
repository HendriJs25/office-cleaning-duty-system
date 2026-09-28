package role

import (
	rolehandler "cleaning/internal/handler/role"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, handler *rolehandler.Handler) {
	router.GET("/roles", handler.GetAll)
}
