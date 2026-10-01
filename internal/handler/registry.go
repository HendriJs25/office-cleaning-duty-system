package handler

import (
	authhandler "cleaning/internal/handler/auth"
	healthhandler "cleaning/internal/handler/health"
	officehandler "cleaning/internal/handler/office"
	rolehandler "cleaning/internal/handler/role"
	userhandler "cleaning/internal/handler/user"
	"cleaning/internal/services"

	customValidator "github.com/go-playground/validator/v10"
)

type Register struct {
	Health *healthhandler.Health
	Auth   *authhandler.Handler
	Role   *rolehandler.Handler
	User   *userhandler.Handler
	Office *officehandler.Handler
}

func NewRegistry(services *services.Registry, validate *customValidator.Validate) *Register {
	return &Register{
		Health: healthhandler.NewHandler(),
		Auth:   authhandler.NewHandler(services.AuthService, validate),
		Role:   rolehandler.NewHandler(services.RoleService),
		User:   userhandler.NewHandler(services.UserService, validate),
		Office: officehandler.NewHandler(services.OfficeService, validate),
	}
}
