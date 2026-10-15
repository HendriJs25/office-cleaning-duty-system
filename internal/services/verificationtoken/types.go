package verificationtoken

import "time"

type CreateRegistrationTokenResult struct {
	RawToken  string
	ExpiresAt time.Time
}
