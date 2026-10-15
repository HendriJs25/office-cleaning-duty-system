package error

import "errors"

var (
	ErrInvalidToken             = errors.New("token invalid")
	ErrExpiredToken             = errors.New("token expired")
	ErrInvalidRegistrationToken = errors.New("invalid registration token")
)

var TokenErrors = []error{
	ErrInvalidToken,
	ErrExpiredToken,
	ErrInvalidRegistrationToken,
}
