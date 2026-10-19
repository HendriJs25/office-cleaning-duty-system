package error

import "errors"

var (
	ErrRoleNotFound        = errors.New("role not found")
	ErrRoleInActive        = errors.New("role inactivate")
	ErrSystemRoleImmutable = errors.New("system role cannot be modified")
	ErrRoleInUse           = errors.New("role is currently assigned to users")
)

var RoleErrors = []error{
	ErrRoleNotFound,
	ErrRoleInActive,
	ErrSystemRoleImmutable,
}
