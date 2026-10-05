package employee

import (
	"cleaning/internal/common/response"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/dto/request"
	employeeservice "cleaning/internal/services/employee"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	customValidator "github.com/go-playground/validator/v10"
)

type Handler struct {
	employeeService employeeservice.Service
	validate        *customValidator.Validate
}

func NewHandler(employeeService employeeservice.Service, validate *customValidator.Validate) *Handler {
	return &Handler{
		employeeService: employeeService,
		validate:        validate,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateEmployeeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	err := h.employeeService.Create(c.Request.Context(), employeeservice.CreateEmployeeInput{
		OfficeID:            req.OfficeID,
		FamilyName:          req.FamilyName,
		GivenName:           req.GivenName,
		FamilyNameKana:      req.FamilyNameKana,
		GivenNameKana:       req.GivenNameKana,
		EmploymentStartDate: req.EmploymentStartDate.TimePtr(),
		EmploymentEndDate:   req.EmploymentEndDate.TimePtr(),
	})

	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrNotFound):
			response.InvalidOption(c, err, "オフィス")
			return
		case errors.Is(err, errConstant.ErrInActive):
			response.InactiveOption(c, err, "オフィス")
			return
		default:
			slog.Error("create employee failed", "error", err)
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "従業員を登録しました。",
		Gin:     c,
	})
}
