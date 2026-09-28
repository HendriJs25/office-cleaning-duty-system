package services

import (
	"cleaning/internal/config"
	"cleaning/internal/repository"
	sessionrepository "cleaning/internal/repository/session"
	userservice "cleaning/internal/services/auth"
	jwtservice "cleaning/internal/services/jwt"
	roleservice "cleaning/internal/services/role"
)

type Registry struct {
	UserService userservice.Service
	JWTService  jwtservice.Service
	RoleService roleservice.Service
}

func NewRegistry(repositories *repository.Registry, sessionRepository sessionrepository.Repository, jwtConfig *config.JWT) *Registry {
	jwtService := jwtservice.NewService(jwtConfig)

	return &Registry{
		UserService: userservice.NewService(repositories.User, sessionRepository, jwtService),
		JWTService:  jwtService,
		RoleService: roleservice.NewService(repositories.Role),
	}
}
