package employee

import (
	errWrap "cleaning/internal/common/error"
	"cleaning/internal/common/response"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/dto/request"
	employeeservice "cleaning/internal/services/employee"
	"errors"
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

	err := h.employeeService.Create(c.Request.Context(), employeeservice.CreateEmployeeInput{
		OfficeID:            req.OfficeID,
		FamilyName:          req.FamilyName,
		GivenName:           req.GivenName,
		FamilyNameKana:      req.FamilyNameKana,
		GivenNameKana:       req.GivenNameKana,
		EmploymentStartDate: req.EmploymentStartDate,
		EmploymentEndDate:   req.EmploymentEndDate,
	})

	if err != nil {

		}
	}
}
