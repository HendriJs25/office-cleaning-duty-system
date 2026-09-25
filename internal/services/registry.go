package services

import (
	"cleaning/internal/config"
	"cleaning/internal/repository"
	sessionrepository "cleaning/internal/repository/session"
	jwtservice "cleaning/internal/services/jwt"
	userservice "cleaning/internal/services/user"
)

type Registry struct {
	UserService userservice.Service
	JWTService  jwtservice.Service
}

func NewRegistry(repositories *repository.Registry, sessionRepository sessionrepository.Repository, jwtConfig *config.JWT) *Registry {
	jwtService := jwtservice.NewService(jwtConfig)

	return &Registry{
		UserService: userservice.NewService(repositories.User, sessionRepository, jwtService),
		JWTService:  jwtService,
	}
}
