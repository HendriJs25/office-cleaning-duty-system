package handler

import (
	healthhandler "cleaning/internal/handler/health"
	userhandler "cleaning/internal/handler/user"
	"cleaning/internal/services"

	customValidator "github.com/go-playground/validator/v10"
)

type Register struct {
	Health *healthhandler.Health
	User   *userhandler.Handler
}

func NewRegistry(services *services.Registry, validate *customValidator.Validate) *Register {
	return &Register{
		Health: healthhandler.NewHandler(),
		User:   userhandler.NewHandler(services.UserService, validate),
	}
}
