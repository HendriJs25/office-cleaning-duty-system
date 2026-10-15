package verificationtoken

import (
	"time"

	"github.com/google/uuid"
)

type CreateRegistrationTokenResult struct {
	RawToken  string
	ExpiresAt time.Time
}

type VerifyRegistrationTokenResult struct {
	EmployeeUUID uuid.UUID
	FullName     string
	FullNameKana *string
	OfficeName   string
}
