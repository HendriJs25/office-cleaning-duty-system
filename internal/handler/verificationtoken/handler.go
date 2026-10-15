package verificationtoken

import (
	"cleaning/internal/common/response"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/dto/request"
	responsedto "cleaning/internal/domain/dto/response"
	"cleaning/internal/logger"
	verificationtokenservice "cleaning/internal/services/verificationtoken"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	customValidator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Handler struct {
	VerificationTokenService verificationtokenservice.Service
	validate                 *customValidator.Validate
}

func NewHandler(verificationTokenService verificationtokenservice.Service, validate *customValidator.Validate) *Handler {
	return &Handler{
		VerificationTokenService: verificationTokenService,
		validate:                 validate,
	}
}

func (h *Handler) CreateRegistrationToken(c *gin.Context) {
	uuidStr := c.Param("uuid")

	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		response.BadRequest(c)
		return
	}

	result, err := h.VerificationTokenService.CreateRegistrationToken(c.Request.Context(), parsedUUID)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrNotFound):
			response.NotFound(c)
			return
		case errors.Is(err, errConstant.ErrEmployeeInActive):
			response.Conflict(c, err)
			return
		case errors.Is(err, errConstant.ErrEmployeeAlreadyAssigned):
			response.Conflict(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("employee_uuid", parsedUUID).WithError(err).Error("error creating registration token")
			response.InternalServerError(c)
			return
		}
	}

	registrationToken := responsedto.RegistrationTokenResponse{
		Token:     result.RawToken,
		ExpiresAt: result.ExpiresAt,
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusCreated,
		Data: registrationToken,
		Gin:  c,
	})
}

func (h *Handler) VerifyRegistrationToken(c *gin.Context) {
	var req request.VerifyRegistrationTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	result, err := h.VerificationTokenService.VerifyRegistrationToken(c.Request.Context(), req.Token)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrInvalidToken),
			errors.Is(err, errConstant.ErrExpiredToken),
			errors.Is(err, errConstant.ErrEmployeeNotFound),
			errors.Is(err, errConstant.ErrEmployeeInActive),
			errors.Is(err, errConstant.ErrEmployeeAlreadyAssigned):
			response.UnprocessableEntity(c, errConstant.ErrInvalidRegistrationToken)
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).Error("error verifying registration token")
			response.InternalServerError(c)
			return
		}
	}

	data := responsedto.VerifyRegistrationTokenResponse{
		EmployeeUUID: result.EmployeeUUID,
		FullName:     result.FullName,
		FullNameKana: result.FullNameKana,
		OfficeName:   result.OfficeName,
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: data,
		Gin:  c,
	})
}
