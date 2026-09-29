package validator

import (
	"log/slog"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	lowerRegex  = regexp.MustCompile(`[a-z]`)
	upperRegex  = regexp.MustCompile(`[A-Z]`)
	numberRegex = regexp.MustCompile(`[0-9]`)
)

func New() *validator.Validate {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if err := v.RegisterValidation("validpassword", validPassword); err != nil {
		slog.Error("custom validator failed to register", "error", err)
		os.Exit(1)
	}

	return v
}

func validPassword(f1 validator.FieldLevel) bool {
	password := f1.Field().String()

	return len(password) >= 8 &&
		lowerRegex.MatchString(password) &&
		upperRegex.MatchString(password) &&
		numberRegex.MatchString(password)
}
