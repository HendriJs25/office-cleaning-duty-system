package office

import (
	errWrap "cleaning/internal/common/error"
	"cleaning/internal/common/response"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/dto/request"
	officeservice "cleaning/internal/services/office"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	customValidator "github.com/go-playground/validator/v10"
)

type Handler struct {
	officeService officeservice.Service
	validate      *customValidator.Validate
}

func NewHandler(officeService officeservice.Service, validate *customValidator.Validate) *Handler {
	return &Handler{officeService: officeService, validate: validate}
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateOfficeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.HTTPResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrBadRequest,
			Gin:  c,
		})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.HTTPResponse(response.ParamHTTPResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: errWrap.ErrValidationResponse(err),
			Err:     err,
			Gin:     c,
		})
		return
	}

	err := h.officeService.Create(c.Request.Context(), officeservice.CreateOfficeInput{
		Code:    req.Code,
		Name:    req.Name,
		Address: req.Address,
	})

	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrAlreadyExists):
			response.HTTPResponse(response.ParamHTTPResponse{
				Code:    http.StatusConflict,
				Message: "このコードが既に登録されています",
				Err:     err,
				Gin:     c,
			})
			return
		default:
			response.HTTPResponse(response.ParamHTTPResponse{
				Code: http.StatusInternalServerError,
				Err:  err,
				Gin:  c,
			})
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "オフィスを作成しました。",
		Gin:     c,
	})
}
