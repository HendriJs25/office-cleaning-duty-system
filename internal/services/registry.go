package services

import (
	"cleaning/internal/config"
	"cleaning/internal/repository"
	permissioncacherepository "cleaning/internal/repository/permissioncache"
	sessionrepository "cleaning/internal/repository/session"
	authservice "cleaning/internal/services/auth"
	employeeservice "cleaning/internal/services/employee"
	jwtservice "cleaning/internal/services/jwt"
	officeservice "cleaning/internal/services/office"
	permissionservice "cleaning/internal/services/permission"
	roleservice "cleaning/internal/services/role"
	userservice "cleaning/internal/services/user"
)

type Registry struct {
	AuthService       authservice.Service
	JWTService        jwtservice.Service
	RoleService       roleservice.Service
	PermissionService permissionservice.Service
	UserService       userservice.Service
	OfficeService     officeservice.Service
	EmployeeService   employeeservice.Service
}

func NewRegistry(repositories *repository.Registry, sessionRepository sessionrepository.Repository, permissionCacheRepository permissioncacherepository.Repository, jwtConfig *config.JWT) *Registry {
	jwtService := jwtservice.NewService(jwtConfig)

	return &Registry{
		AuthService:       authservice.NewService(repositories.User, sessionRepository, jwtService),
		JWTService:        jwtService,
		RoleService:       roleservice.NewService(repositories.Role),
		PermissionService: permissionservice.NewService(repositories.Permission, permissionCacheRepository),
		UserService:       userservice.NewService(repositories.User, repositories.Role, sessionRepository),
		OfficeService:     officeservice.NewService(repositories.Office),
		EmployeeService:   employeeservice.NewService(repositories.Employee, repositories.Office),
	}
}
