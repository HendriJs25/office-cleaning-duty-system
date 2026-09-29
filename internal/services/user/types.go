package user

import (
	"time"

	"github.com/google/uuid"
)

type GetUserResult struct {
	UserID      int64
	UUID        uuid.UUID
	UserName    string
	Email       string
	IsActive    bool
	LastLoginAt *time.Time
	RoleName    string
}

type CreateUserInput struct {
	RoleID     int64
	EmployeeID *int64
	UserName   string
	Email      string
	Password   string
}
