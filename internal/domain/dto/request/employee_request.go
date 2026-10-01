package request

import "time"

type CreateEmployeeRequest struct {
	OfficeID            int64      `json:"office_id" validate:"required"`
	FamilyName          string     `json:"family_name" validate:"required,notblank,max=100"`
	GivenName           string     `json:"given_name" validate:"required,notblank,max=100"`
	FamilyNameKana      *string    `json:"family_name_kana" validate:"omitempty,notblank,max=100"`
	GivenNameKana       *string    `json:"given_name_kana" validate:"omitempty,notblank,max=100"`
	EmploymentStartDate *time.Time `json:"employment_start_date" validate:"required_with=EmploymentEndDate"`
	EmploymentEndDate   *time.Time `json:"employment_end_date" validate:"omitempty,gtfield=EmploymentStartDate"`
}
