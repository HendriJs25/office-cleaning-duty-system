package error

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationResponse struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

var ErrValidator = map[string]string{
	"required":      "%sは必須です",
	"email":         "%sの形式は正しくありません",
	"max":           "%sは%s文字以内で入力してください",
	"validpassword": "%sは8文字以上で、大文字・小文字・数字をそれぞれ1文字以上含めてください",
	"eqfield":       "%sがパスワードと一致していません",
	"notblank":      "%sは空白のみでは登録できません",
}

var FieldNameJa = map[string]string{
	"email":            "メールアドレス",
	"password":         "パスワード",
	"user_name":        "ユーザー名",
	"password_confirm": "パスワード確認",
}

func ErrValidationResponse(err error) (validationResponse []ValidationResponse) {
	var fieldErrors validator.ValidationErrors

	if errors.As(err, &fieldErrors) {
		for _, err := range fieldErrors {
			fieldJa := getFieldNameJa(err.Field())
			errValidator, ok := ErrValidator[err.Tag()]
			if ok {
				count := strings.Count(errValidator, "%s")
				if count == 1 {
					validationResponse = append(validationResponse, ValidationResponse{
						Field:   err.Field(),
						Message: fmt.Sprintf(errValidator, fieldJa),
					})
				} else {
					validationResponse = append(validationResponse, ValidationResponse{
						Field:   err.Field(),
						Message: fmt.Sprintf(errValidator, fieldJa, err.Param()),
					})
				}
			} else {
				validationResponse = append(validationResponse, ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%sの入力内容が正しくありません", fieldJa),
				})
			}
		}
	}
	return validationResponse
}

func getFieldNameJa(field string) string {
	if name, ok := FieldNameJa[field]; ok {
		return name
	}
	return field
}
