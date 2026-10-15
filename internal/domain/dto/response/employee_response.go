package response

import (
	"time"

	"github.com/google/uuid"
)

type GetActiveEmployeesResponse struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
}

type GetEmployeesResponse struct {
	ID                  int64      `json:"id"`
	UUID                uuid.UUID  `json:"uuid"`
	OfficeName          string     `json:"office_name"`
	FullName            string     `json:"full_name"`
	IsActive            bool       `json:"is_active"`
	EmploymentStartDate *time.Time `json:"employment_start_date"`
}

type GetEmployeeDetailResponse struct {
	ID                  int64          `json:"id"`
	UUID                uuid.UUID      `json:"uuid"`
	OfficeName          string         `json:"office_name"`
	FullName            string         `json:"full_name"`
	FullNameKana        string         `json:"full_name_kana"`
	IsActive            bool           `json:"is_active"`
	EmploymentStartDate *time.Time     `json:"employment_start_date"`
	EmploymentEndDate   *time.Time     `json:"employment_end_date"`
	Token               *TokenResponse `json:"token"`
}

type TokenResponse struct {
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
}
