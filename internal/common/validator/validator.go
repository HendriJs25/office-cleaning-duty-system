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
	lowerRegex           = regexp.MustCompile(`[a-z]`)
	upperRegex           = regexp.MustCompile(`[A-Z]`)
	numberRegex          = regexp.MustCompile(`[0-9]`)
	lowercaseHyphenRegex = regexp.MustCompile(`^[a-z-]+$`)
	katakanaRegex        = regexp.MustCompile(`^[ァ-ヶー]+$`)
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

	registerCustomValidators(v)

	return v
}

func registerCustomValidators(v *validator.Validate) {
	validations := map[string]validator.Func{
		"validpassword":    validPassword,
		"notblank":         notBlank,
		"lowercase_hyphen": lowercaseHyphen,
		"katakana":         validateKatakana,
	}

	for tag, fn := range validations {
		if err := v.RegisterValidation(tag, fn); err != nil {
			slog.Error("custom validator failed to register",
				"tag", tag,
				"error", err)
			os.Exit(1)
		}
	}
}

func validPassword(f1 validator.FieldLevel) bool {
	password := f1.Field().String()

	return len(password) >= 8 &&
		lowerRegex.MatchString(password) &&
		upperRegex.MatchString(password) &&
		numberRegex.MatchString(password)
}

func notBlank(f1 validator.FieldLevel) bool {
	return strings.TrimSpace(f1.Field().String()) != ""
}

func lowercaseHyphen(f1 validator.FieldLevel) bool {
	return lowercaseHyphenRegex.MatchString(f1.Field().String())
}

func validateKatakana(f1 validator.FieldLevel) bool {
	return katakanaRegex.MatchString(f1.Field().String())
}
