package response

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	UserID      int64      `json:"user_id"`
	UUID        uuid.UUID  `json:"uuid"`
	Username    string     `json:"user_name"`
	Email       string     `json:"email"`
	IsActive    *bool      `json:"is_active,omitempty"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	RoleName    string     `json:"role_name"`
}

type LoginResponse struct {
	User *UserResponse `json:"auth"`
}
