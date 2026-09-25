package user

import (
	errWrap "cleaning/internal/common/error"
	"cleaning/internal/common/response"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/dto/request"
	responsedto "cleaning/internal/domain/dto/response"
	userservice "cleaning/internal/services/user"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	customValidator "github.com/go-playground/validator/v10"
)

type Handler struct {
	userService userservice.Service
	validate    *customValidator.Validate
}

func NewHandler(userService userservice.Service, validate *customValidator.Validate) *Handler {
	return &Handler{
		userService: userService,
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

	result, err := h.userService.Login(c.Request.Context(), userservice.LoginInput{
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

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: responsedto.LoginResponse{
			User: &responsedto.UserResponse{
				UUID:     result.User.UUID,
				Email:    result.User.Email,
				Username: result.User.UserName,
				RoleName: result.User.RoleName,
			},
		},
		Token: &result.Token,
		Gin:   c,
	})
}
