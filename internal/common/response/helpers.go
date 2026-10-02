package response

import (
	errWrap "cleaning/internal/common/error"
	errConstant "cleaning/internal/constants/error"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequest(c *gin.Context) {
	HTTPResponse(ParamHTTPResponse{
		Code: http.StatusBadRequest,
		Err:  errConstant.ErrBadRequest,
		Gin:  c,
	})
}

func ValidationError(c *gin.Context, err error) {
	HTTPResponse(ParamHTTPResponse{
		Code:    http.StatusUnprocessableEntity,
		Message: errWrap.ErrValidationResponse(err),
		Err:     err,
		Gin:     c,
	})
}

func Unauthorized(c *gin.Context) {
	HTTPResponse(ParamHTTPResponse{
		Code: http.StatusUnauthorized,
		Err:  errConstant.ErrInvalidEmailOrPassword,
		Gin:  c,
	})
}

func Forbidden(c *gin.Context, err error) {
	HTTPResponse(ParamHTTPResponse{
		Code: http.StatusForbidden,
		Err:  err,
		Gin:  c,
	})
}

func InternalServerError(c *gin.Context, err error) {
	HTTPResponse(ParamHTTPResponse{
		Code: http.StatusInternalServerError,
		Err:  err,
		Gin:  c,
	})
}

func Conflict(c *gin.Context, err error, field string) {
	HTTPResponse(ParamHTTPResponse{
		Code:    http.StatusConflict,
		Message: fmt.Sprintf("%sが既に登録されています", field),
		Err:     err,
		Gin:     c,
	})
}

func Empty(c *gin.Context, resource string) {
	HTTPResponse(ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("%sデータがありません", resource),
		Data:    []any{},
		Gin:     c,
	})
}
