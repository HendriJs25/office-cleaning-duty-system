package routes

import (
	"cleaning/internal/handler"
	"cleaning/internal/middleware"
	authroutes "cleaning/internal/routes/auth"
	healthroutes "cleaning/internal/routes/health"

	"github.com/gin-gonic/gin"
)

type Register struct {
	router         *gin.RouterGroup
	handlers       *handler.Register
	authentication *middleware.Authentication
}

func NewRegistry(router *gin.RouterGroup, handlers *handler.Register, authentication *middleware.Authentication) *Register {
	return &Register{
		router:         router,
		handlers:       handlers,
		authentication: authentication,
	}
}

func (r *Register) Register() {
	healthroutes.Register(r.router, r.handlers.Health)
	authroutes.Register(r.router, r.handlers.Auth, r.handlers.Role, r.authentication)
}
