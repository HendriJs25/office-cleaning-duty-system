package error

import (
	"errors"
	"slices"
)

func ErrMapping(err error) bool {
	allErrors := make([]error, 0)
	allErrors = slices.Concat(GeneralErrors, UserErrors)

	for _, item := range allErrors {
		if errors.Is(err, item) {
			return true
		}
	}
	return false
}
