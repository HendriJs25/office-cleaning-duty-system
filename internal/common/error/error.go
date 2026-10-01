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
	"required":         "%sは必須です",
	"email":            "%sの形式は正しくありません",
	"max":              "%sは%s文字以内で入力してください",
	"min":              "%sは%s文字以上で入力してください",
	"validpassword":    "%sは8文字以上で、大文字・小文字・数字をそれぞれ1文字以上含めてください",
	"eqfield":          "%sがパスワードと一致していません",
	"notblank":         "%sは空白のみでは登録できません",
	"lowercase_hyphen": "%sは半角小文字とハイフン（-）のみ使用できます",
	"required_with":    "%sは%sが入力されている場合、必須です",
	"gtfield":          "%sは%sより後の日時を指定してください",
}

var FieldNameJa = map[string]string{
	"email":                 "メールアドレス",
	"password":              "パスワード",
	"user_name":             "ユーザー名",
	"password_confirm":      "パスワード確認",
	"code":                  "コード",
	"name":                  "名前",
	"address":               "住所",
	"office_id":             "オフィス",
	"family_name":           "姓",
	"given_name":            "名",
	"family_name_kana":      "姓（カナ）",
	"given_name_kana":       "名（カナ）",
	"employment_start_date": "勤務開始日",
	"employment_end_date":   "勤務終了日",
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
