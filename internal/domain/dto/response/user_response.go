package response

import "github.com/google/uuid"

type UserResponse struct {
	UUID     uuid.UUID `json:"uuid"`
	Email    string    `json:"email"`
	Username string    `json:"user_name"`
	RoleName string    `json:"role_name"`
}

type LoginResponse struct {
	User *UserResponse `json:"user"`
}
