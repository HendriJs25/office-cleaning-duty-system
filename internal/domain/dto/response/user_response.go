package response

import (
	"time"

	"github.com/google/uuid"
)

type UserInfoResponse struct {
	ID          int64      `json:"id"`
	UUID        uuid.UUID  `json:"uuid"`
	UserName    string     `json:"user_name"`
	Email       string     `json:"email"`
	IsActive    *bool      `json:"is_active,omitempty"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	RoleName    string     `json:"role_name"`
}

type EmployeeInfoResponse struct {
	ID                  int64      `json:"id"`
	UUID                uuid.UUID  `json:"uuid"`
	OfficeName          string     `json:"office_name"`
	FullName            string     `json:"full_name"`
	EmploymentStartDate *time.Time `json:"employment_start_date,omitempty"`
}

type UserDetailResponse struct {
	User     UserInfoResponse      `json:"user"`
	Employee *EmployeeInfoResponse `json:"employee,omitempty"`
}

type EmployeeProfileResponse struct {
	OfficeName          string     `json:"office_name"`
	FullName            string     `json:"full_name"`
	FullNameKana        *string    `json:"full_name_kana"`
	EmploymentStartDate *time.Time `json:"employment_start_date,omitempty"`
}

type UserProfileResponse struct {
	UserName string                   `json:"user_name"`
	Email    string                   `json:"email"`
	RoleName string                   `json:"role_name"`
	Employee *EmployeeProfileResponse `json:"employee,omitempty"`
}
type LoginResponse struct {
	User *UserInfoResponse `json:"auth"`
}
