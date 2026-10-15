package response

import "github.com/google/uuid"

type VerifyRegistrationTokenResponse struct {
	EmployeeUUID uuid.UUID `json:"employee_uuid"`
	FullName     string    `json:"full_name"`
	FullNameKana *string   `json:"full_name_kana"`
	OfficeName   string    `json:"office_name"`
}
