package error

import "errors"

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrPasswordIncorrect      = errors.New("password incorrect")
	ErrAccountIsDeactivated   = errors.New("account is deactivated")
)

var UserErrors = []error{
	ErrInvalidEmailOrPassword,
	ErrPasswordIncorrect,
	ErrAccountIsDeactivated,
}
