package office

import (
	errWrap "cleaning/internal/common/error"
	"cleaning/internal/common/response"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/dto/request"
	responsedto "cleaning/internal/domain/dto/response"
	officeservice "cleaning/internal/services/office"
	"errors"
	"log/slog"
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
			slog.Error("create office failed", "error", err)
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

func (h *Handler) GetAllActiveOffices(c *gin.Context) {
	result, err := h.officeService.GetAllActiveOffices(c.Request.Context())
	if err != nil {
		slog.Error("get all office failed", "error", err)
		response.HTTPResponse(response.ParamHTTPResponse{
			Code: http.StatusInternalServerError,
			Err:  err,
			Gin:  c,
		})
		return
	}

	if len(result) == 0 {
		response.HTTPResponse(response.ParamHTTPResponse{
			Code:    http.StatusOK,
			Message: "オフィスがありません",
			Data:    nil,
			Gin:     c,
		})
		return
	}

	var offices []responsedto.GetActiveOfficesResponse
	for _, office := range result {
		offices = append(offices, responsedto.GetActiveOfficesResponse{
			ID:   office.ID,
			Name: office.Name,
		})
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: offices,
		Gin:  c,
	})

}
