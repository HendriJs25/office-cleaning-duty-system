package error

import "errors"

var (
	ErrOfficeInUse = errors.New("office is in use")
)

var OfficeErrors = []error{
	ErrOfficeInUse,
}
