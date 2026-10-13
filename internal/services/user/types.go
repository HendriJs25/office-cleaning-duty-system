package user

import (
	"time"

	"github.com/google/uuid"
)

type GetUserInfoResult struct {
	ID          int64
	UUID        uuid.UUID
	UserName    string
	Email       string
	IsActive    bool
	LastLoginAt *time.Time
	RoleName    string
}

type GetEmployeeInfoResult struct {
	ID                  int64
	UUID                uuid.UUID
	OfficeName          string
	FullName            string
	EmploymentStartDate *time.Time
}

type GetUserDetailResult struct {
	User     GetUserInfoResult
	Employee *GetEmployeeInfoResult
}
type GetEmployeeProfileResult struct {
	OfficeName          string
	FullName            string
	FullNameKana        *string
	EmploymentStartDate *time.Time
}

type GetUserProfileResult struct {
	UserName string
	Email    string
	RoleName string
	Employee *GetEmployeeProfileResult
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
