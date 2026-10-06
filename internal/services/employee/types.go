package employee

import "time"

type CreateEmployeeInput struct {
	OfficeID            int64
	FamilyName          string
	GivenName           string
	FamilyNameKana      *string
	GivenNameKana       *string
	EmploymentStartDate *time.Time
	EmploymentEndDate   *time.Time
}

type GetActiveEmployeesResult struct {
	ID       int64
	FullName string
}
