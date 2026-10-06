package office

import (
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
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
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
			response.ConflictDuplicate(c, err, "コード")
			return
		default:
			slog.Error("create office failed", "error", err)
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "オフィスを登録しました",
		Gin:     c,
	})
}

func (h *Handler) GetAllActiveOffices(c *gin.Context) {
	offices, err := h.officeService.GetAllActiveOffices(c.Request.Context())
	if err != nil {
		slog.Error("get all office failed", "error", err)
		response.InternalServerError(c)
		return
	}

	if len(offices) == 0 {
		response.Empty(c, "オフィス")
		return
	}

	var result []responsedto.GetActiveOfficesResponse
	for _, office := range offices {
		result = append(result, responsedto.GetActiveOfficesResponse{
			ID:   office.ID,
			Name: office.Name,
		})
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}
