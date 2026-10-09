package error

import "errors"

var (
	ErrEmployeeNotFound        = errors.New("employee not found")
	ErrEmployeeInActive        = errors.New("employee inactivate")
	ErrEmployeeAlreadyAssigned = errors.New("employee already assigned")
)

var EmployeeErrors = []error{
	ErrEmployeeNotFound,
	ErrEmployeeInActive,
	ErrEmployeeAlreadyAssigned,
}
