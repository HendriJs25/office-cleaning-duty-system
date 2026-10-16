package handler

import (
	authhandler "cleaning/internal/handler/auth"
	employeehandler "cleaning/internal/handler/employee"
	healthhandler "cleaning/internal/handler/health"
	officehandler "cleaning/internal/handler/office"
	permissionhandler "cleaning/internal/handler/permission"
	rolehandler "cleaning/internal/handler/role"
	userhandler "cleaning/internal/handler/user"
	verificationtokenhandler "cleaning/internal/handler/verificationtoken"
	"cleaning/internal/services"

	customValidator "github.com/go-playground/validator/v10"
)

type Register struct {
	Health            *healthhandler.Health
	Auth              *authhandler.Handler
	Role              *rolehandler.Handler
	User              *userhandler.Handler
	Office            *officehandler.Handler
	Employee          *employeehandler.Handler
	VerificationToken *verificationtokenhandler.Handler
	Permission        *permissionhandler.Handler
}

func NewRegistry(services *services.Registry, validate *customValidator.Validate) *Register {
	return &Register{
		Health:            healthhandler.NewHandler(),
		Auth:              authhandler.NewHandler(services.AuthService, validate),
		Role:              rolehandler.NewHandler(services.RoleService, validate),
		User:              userhandler.NewHandler(services.UserService, validate),
		Office:            officehandler.NewHandler(services.OfficeService, validate),
		Employee:          employeehandler.NewHandler(services.EmployeeService, validate),
		VerificationToken: verificationtokenhandler.NewHandler(services.VerificationTokenService, validate),
		Permission:        permissionhandler.NewHandler(services.PermissionService),
	}
}
