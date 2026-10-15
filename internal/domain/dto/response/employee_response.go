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
