package name

import (
	"fmt"
	"strings"
)

func BuildFullName(familyName, givenName string) string {
	return fmt.Sprintf("%s %s", familyName, givenName)
}

func BuildFullNameKana(familyNameKana, givenNameKana *string) *string {
	var names []string

	if familyNameKana != nil {
		names = append(names, strings.TrimSpace(*familyNameKana))
	}

	if givenNameKana != nil {
		names = append(names, strings.TrimSpace(*givenNameKana))
	}

	if len(names) == 0 {
		return nil
	}

	fullName := strings.Join(names, " ")
	return &fullName
}
