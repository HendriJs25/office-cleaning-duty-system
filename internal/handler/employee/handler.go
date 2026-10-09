package employee

import (
	"cleaning/internal/common/response"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/dto/request"
	responsedto "cleaning/internal/domain/dto/response"
	"cleaning/internal/logger"
	employeeservice "cleaning/internal/services/employee"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	customValidator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
			logger.WithContext(c.Request.Context()).WithError(err).Error("Create employee failed")
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

func (h *Handler) GetAllAssignableEmployees(c *gin.Context) {
	var userUUID *uuid.UUID

	if uuidStr := c.Query("user_uuid"); uuidStr != "" {
		parsedUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			response.BadRequest(c)
			return
		}

		userUUID = &parsedUUID
	}

	employees, err := h.employeeService.GetAllAssignableEmployees(c.Request.Context(), userUUID)
	if err != nil {
		logger.WithContext(c.Request.Context()).WithError(err).Error("Get all active employees failed")
		response.InternalServerError(c)
		return
	}

	if len(employees) == 0 {
		response.Empty(c, "従業員")
		return
	}

	var result []responsedto.GetActiveEmployeesResponse
	for _, employee := range employees {
		result = append(result, responsedto.GetActiveEmployeesResponse{
			ID:       employee.ID,
			FullName: employee.FullName,
		})
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}
