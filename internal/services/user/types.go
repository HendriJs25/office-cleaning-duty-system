package user

import (
	"time"

	"github.com/google/uuid"
)

type GetUserResult struct {
	ID          int64
	UUID        uuid.UUID
	UserName    string
	Email       string
	IsActive    bool
	LastLoginAt *time.Time
	RoleName    string
}

type EmployeeInfo struct {
	ID                  int64
	UUID                uuid.UUID
	OfficeName          string
	FullName            string
	EmploymentStartDate *time.Time
}

type GetUserDetailResult struct {
	User     GetUserResult
	Employee *EmployeeInfo
}

type CreateUserInput struct {
	RoleID     int64
	EmployeeID *int64
	UserName   string
	Email      string
	Password   string
}

type UpdateUserByAdminInput struct {
	RoleID     int64
	EmployeeID *int64
	IsActive   bool
}
