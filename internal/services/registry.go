package services

import (
	"cleaning/internal/config"
	"cleaning/internal/repository"
	permissioncacherepository "cleaning/internal/repository/permissioncache"
	sessionrepository "cleaning/internal/repository/session"
	authservice "cleaning/internal/services/auth"
	jwtservice "cleaning/internal/services/jwt"
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
}

func NewRegistry(repositories *repository.Registry, sessionRepository sessionrepository.Repository, permissionCacheRepository permissioncacherepository.Repository, jwtConfig *config.JWT) *Registry {
	jwtService := jwtservice.NewService(jwtConfig)

	return &Registry{
		AuthService:       authservice.NewService(repositories.User, sessionRepository, jwtService),
		JWTService:        jwtService,
		RoleService:       roleservice.NewService(repositories.Role),
		PermissionService: permissionservice.NewService(repositories.Permission, permissionCacheRepository),
		UserService:       userservice.NewService(repositories.User, repositories.Role),
	}
}
