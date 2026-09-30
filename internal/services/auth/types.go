package auth

import (
	"cleaning/internal/services/jwt"

	"github.com/google/uuid"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken *jwt.AccessToken
	User        AuthenticatedUser
}

type AuthenticateInput struct {
	Email    string
	Password string
}

type AuthenticatedUser struct {
	ID       int64
	UUID     uuid.UUID
	Email    string
	UserName string
	RoleID   int64
	RoleName string
}
