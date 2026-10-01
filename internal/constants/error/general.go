package error

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrTooManyRequests     = errors.New("too many requests")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrBadRequest          = errors.New("bad request")
	ErrNotFound            = errors.New("not found")
	ErrAlreadyExists       = errors.New("already exists")
	ErrInActive            = errors.New("inactive")

	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token is expired")
)

var GeneralErrors = []error{
	ErrInternalServerError,
	ErrTooManyRequests,
	ErrUnauthorized,
	ErrForbidden,
	ErrInvalidArgument,
	ErrBadRequest,
	ErrNotFound,
	ErrAlreadyExists,
	ErrInActive,
	ErrInvalidToken,
	ErrTokenExpired,
}
