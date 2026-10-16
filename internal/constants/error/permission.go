package error

import "errors"

var (
	ErrInvalidPermission = errors.New("invalid permission")
)

var PermissionErrors = []error{
	ErrInvalidPermission,
}
