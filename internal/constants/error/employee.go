package error

import "errors"

var (
	ErrEmployeeNotFound = errors.New("employee not found")
	ErrEmployeeInActive = errors.New("employee inactivate")
)

var EmployeeErrors = []error{
	ErrEmployeeNotFound,
	ErrEmployeeInActive,
}
