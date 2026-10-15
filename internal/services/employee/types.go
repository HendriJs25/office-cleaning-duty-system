package employee

import (
	"time"

	"github.com/google/uuid"
)

type CreateEmployeeInput struct {
	OfficeID            int64
	FamilyName          string
	GivenName           string
	FamilyNameKana      *string
	GivenNameKana       *string
	EmploymentStartDate *time.Time
	EmploymentEndDate   *time.Time
}

type GetAssignableEmployeesResult struct {
	ID       int64
	FullName string
}

type GetEmployeeResult struct {
	ID                  int64
	UUID                uuid.UUID
	OfficeName          string
	FullName            string
	IsActive            bool
	EmploymentStartDate *time.Time
}

type InvitationToken struct {
	Token     string
	ExpiresAt time.Time
}

type UpdateEmployeeByAdminInput struct {
	OfficeID            int64
	FamilyName          string
	GivenName           string
	FamilyNameKana      *string
	GivenNameKana       *string
	IsActive            bool
	EmploymentStartDate *time.Time
	EmploymentEndDate   *time.Time
}
