package error

import "errors"

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrPasswordIncorrect      = errors.New("password incorrect")
	ErrAccountIsDeactivated   = errors.New("account is deactivated")
	ErrAlreadyDeactivated     = errors.New("account is already deactivated")
	ErrAlreadyActivated       = errors.New("account is already activated")
)

var UserErrors = []error{
	ErrInvalidEmailOrPassword,
	ErrPasswordIncorrect,
	ErrAccountIsDeactivated,
	ErrAlreadyDeactivated,
	ErrAlreadyActivated,
}
