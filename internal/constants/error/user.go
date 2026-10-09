package error

import "errors"

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrPasswordIncorrect      = errors.New("password incorrect")
	ErrAccountIsDeactivated   = errors.New("account is deactivated")
	ErrCannotUpdateSelf       = errors.New("cannot update self")
)

var UserErrors = []error{
	ErrInvalidEmailOrPassword,
	ErrPasswordIncorrect,
	ErrAccountIsDeactivated,
	ErrCannotUpdateSelf,
}
