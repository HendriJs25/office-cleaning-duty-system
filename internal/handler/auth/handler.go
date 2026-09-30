package auth

import (
	"cleaning/internal/common/cookie"
	errWrap "cleaning/internal/common/error"
	"cleaning/internal/common/response"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/dto/request"
	responsedto "cleaning/internal/domain/dto/response"
	"cleaning/internal/middleware"
	authservice "cleaning/internal/services/auth"
	"errors"
	"log/slog"
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

	result, err := h.authService.Login(c.Request.Context(), authservice.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrNotFound) || errors.Is(err, errConstant.ErrPasswordIncorrect):
			response.HTTPResponse(response.ParamHTTPResponse{
				Code: http.StatusUnauthorized,
				Err:  errConstant.ErrInvalidEmailOrPassword,
				Gin:  c,
			})
			return
		case errors.Is(err, errConstant.ErrAccountIsDeactivated):
			response.HTTPResponse(response.ParamHTTPResponse{
				Code: http.StatusForbidden,
				Err:  err,
				Gin:  c,
			})
			return
		default:
			slog.Error("login failed", "email", req.Email, "error", err)
			response.HTTPResponse(response.ParamHTTPResponse{
				Code: http.StatusInternalServerError,
				Err:  err,
				Gin:  c,
			})
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
		slog.Error("logout failed", "error", err)
		response.HTTPResponse(response.ParamHTTPResponse{
			Code: http.StatusInternalServerError,
			Err:  errConstant.ErrInternalServerError,
			Gin:  c,
		})
	}

	cookie.ClearAccessToken(c)

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "ログアウトしました。",
		Gin:     c,
	})
}
