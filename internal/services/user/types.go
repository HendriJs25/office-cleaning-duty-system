package user

import "github.com/google/uuid"

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	Token string
	User  AuthenticatedUser
}

type AuthenticateInput struct {
	Email    string
	Password string
}

type AuthenticatedUser struct {
	UUID     uuid.UUID
	Email    string
	RoleCode string
	RoleName string
}
