package response

import (
	"cleaning/internal/constants"
	errConstant "cleaning/internal/constants/error"
	"errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string `json:"status"`
	Message any    `json:"message,omitempty"`
	Data    any    `json:"data"`
}

type ParamHTTPResponse struct {
	Code    int
	Err     error
	Message any
	Gin     *gin.Context
	Data    any
}

func HTTPResponse(param ParamHTTPResponse) {
	if param.Err == nil {
		param.Gin.JSON(param.Code, Response{
			Status:  constants.Success,
			Message: param.Message,
			Data:    param.Data,
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
	case errors.Is(err, errConstant.ErrAlreadyDeactivated):
		return "このアカウントは既に無効化されています。"
	case errors.Is(err, errConstant.ErrAlreadyActivated):
		return "このアカウントは既に効化しています。"

	case errors.Is(err, errConstant.ErrNotFound):
		return "該当するデータがありません。"
	case errors.Is(err, errConstant.ErrBadRequest):
		return "入力内容に誤りがあります。"
	default:
		return "サーバー内部でエラーが発生しました。"
	}
}
