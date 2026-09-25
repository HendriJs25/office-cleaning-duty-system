package response

import (
	"cleaning/internal/constants"
	errConstant "cleaning/internal/constants/error"
	"errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string  `json:"status"`
	Message any     `json:"message,omitempty"`
	Data    any     `json:"data"`
	Token   *string `json:"token,omitempty"`
}

type ParamHTTPResponse struct {
	Code    int
	Err     error
	Message any
	Gin     *gin.Context
	Data    any
	Token   *string
}

func HTTPResponse(param ParamHTTPResponse) {
	if param.Err == nil {
		param.Gin.JSON(param.Code, Response{
			Status:  constants.Success,
			Message: param.Message,
			Data:    param.Data,
			Token:   param.Token,
		})
		return
	}
	var message any
	if param.Message != nil {
		message = param.Message
	} else {
		message = messageJa(param.Err)
	}

	param.Gin.JSON(param.Code, Response{
		Status:  constants.Error,
		Message: message,
		Data:    param.Data,
	})
	return
}

func messageJa(err error) string {
	switch {
	case errors.Is(err, errConstant.ErrInvalidEmailOrPassword):
		return "メールアドレスまたはパスワードが正しくありません。"
	case errors.Is(err, errConstant.ErrAccountIsDeactivated):
		return "このアカウントは無効化されています。管理者にお問い合わせください。"
	default:
		return "サーバー内部でエラーが発生しました。"
	}
}
