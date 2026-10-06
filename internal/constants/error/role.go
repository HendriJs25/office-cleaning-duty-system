package error

import "errors"

var (
	ErrRoleNotFound = errors.New("employee not found")
	ErrRoleInActive = errors.New("employee inactivate")
)

var RoleErrors = []error{
	ErrRoleNotFound,
	ErrRoleInActive,
}
