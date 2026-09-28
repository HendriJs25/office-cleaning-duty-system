package handler

import (
	userhandler "cleaning/internal/handler/auth"
	healthhandler "cleaning/internal/handler/health"
	rolehandler "cleaning/internal/handler/role"
	"cleaning/internal/services"

	customValidator "github.com/go-playground/validator/v10"
)

type Register struct {
	Health *healthhandler.Health
	Auth   *userhandler.Handler
	Role   *rolehandler.Handler
}

func NewRegistry(services *services.Registry, validate *customValidator.Validate) *Register {
	return &Register{
		Health: healthhandler.NewHandler(),
		Auth:   userhandler.NewHandler(services.UserService, validate),
		Role:   rolehandler.NewHandler(services.RoleService),
	}
}
