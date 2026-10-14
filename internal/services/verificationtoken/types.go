package verificationtoken

import "time"

type CreateRegistrationTokenResult struct {
	rawToken  string
	ExpiresAt time.Time
}
