package auth

import (
	"cleaning/internal/common/cookie"
	"cleaning/internal/common/response"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/dto/request"
	responsedto "cleaning/internal/domain/dto/response"
	"cleaning/internal/logger"
	"cleaning/internal/middleware"
	authservice "cleaning/internal/services/auth"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	customValidator "github.com/go-playground/validator/v10"
)

type Handler struct {
	authService authservice.Service
	validate    *customValidator.Validate
}

func NewHandler(userService authservice.Service, validate *customValidator.Validate) *Handler {
	return &Handler{
		authService: userService,
		validate:    validate,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	result, err := h.authService.Login(c.Request.Context(), authservice.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrNotFound) || errors.Is(err, errConstant.ErrPasswordIncorrect):
			response.Unauthorized(c, errConstant.ErrInvalidEmailOrPassword)
			return
		case errors.Is(err, errConstant.ErrAccountIsDeactivated):
			response.Forbidden(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).Error("Login failed")
			response.InternalServerError(c)
			return
		}
	}

	cookieCreatedAt := time.Now().UTC()
	cookieTTL := result.AccessToken.ExpiresAt.Sub(cookieCreatedAt)
	cookie.SetAccessToken(c, result.AccessToken.Value, cookieTTL)

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: responsedto.LoginResponse{
			User: &responsedto.UserResponse{
				ID:       result.User.ID,
				UUID:     result.User.UUID,
				Email:    result.User.Email,
				Username: result.User.UserName,
				RoleName: result.User.RoleName,
			},
		},
		Gin: c,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	identity, ok := middleware.RequireIdentity(c)
	if !ok {
		return
	}

	err := h.authService.Logout(c.Request.Context(), identity.Token)
	if err != nil {
		logger.WithContext(c.Request.Context()).WithError(err).Error("Logout failed")
		response.InternalServerError(c)
	}

	cookie.ClearAccessToken(c)

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "ログアウトしました",
		Gin:     c,
	})
}
