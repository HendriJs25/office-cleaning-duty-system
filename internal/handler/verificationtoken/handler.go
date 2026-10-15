package verificationtoken

import (
	"cleaning/internal/common/response"
	errConstant "cleaning/internal/constants/error"
	responsedto "cleaning/internal/domain/dto/response"
	"cleaning/internal/logger"
	verificationtokenservice "cleaning/internal/services/verificationtoken"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	VerificationTokenService verificationtokenservice.Service
}

func NewHandler(verificationTokenService verificationtokenservice.Service) *Handler {
	return &Handler{
		VerificationTokenService: verificationTokenService,
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
