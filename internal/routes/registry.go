package routes

import (
	"cleaning/internal/handler"
	healthroutes "cleaning/internal/routes/health"

	"github.com/gin-gonic/gin"
)

type Register struct {
	router   *gin.RouterGroup
	handlers *handler.Register
}

func NewRegistry(router *gin.RouterGroup, handlers *handler.Register) *Register {
	return &Register{
		router:   router,
		handlers: handlers,
	}
}

func (r *Register) Register() {
	healthroutes.Register(r.router, r.handlers.Health)
}
