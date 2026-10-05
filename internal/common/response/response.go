package response

import (
	"cleaning/internal/constants"

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
