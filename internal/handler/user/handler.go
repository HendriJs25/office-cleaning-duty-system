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
	UserService userservice.Service
	validate    *customValidator.Validate
}

func NewHandler(userService userservice.Service, validate *customValidator.Validate) *Handler {
	return &Handler{
		UserService: userService,
		validate:    validate,
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	result, err := h.UserService.GetAll(c.Request.Context())
	if err != nil {
		slog.Error("get all users failed", "error", err)
		response.HTTPResponse(response.ParamHTTPResponse{
			Code: http.StatusInternalServerError,
			Err:  err,
			Data: nil,
			Gin:  c,
		})
		return
	}

	var users []responsedto.UserResponse
	for _, user := range result {
		users = append(users, responsedto.UserResponse{
			UserID:      user.UserID,
			UUID:        user.UUID,
			Username:    user.UserName,
			Email:       user.Email,
			IsActive:    &user.IsActive,
			LastLoginAt: user.LastLoginAt,
			RoleName:    user.RoleName,
		})
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: users,
		Gin:  c,
	})
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateUserRequest

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

	err := h.UserService.Create(c.Request.Context(), userservice.CreateUserInput{
		RoleID:     req.RoleID,
		EmployeeID: req.EmployeeID,
		UserName:   req.UserName,
		Email:      req.Email,
		Password:   req.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrNotFound):
			response.HTTPResponse(response.ParamHTTPResponse{
				Code:    http.StatusOK,
				Message: "ロールがありません。",
				Gin:     c,
			})
			return
		case errors.Is(err, errConstant.ErrAlreadyExists):
			response.HTTPResponse(response.ParamHTTPResponse{
				Code:    http.StatusOK,
				Message: "このメールアドレスは既に登録されています。",
				Gin:     c,
			})
			return
		default:
			slog.Error("create user failed", "error", err)
			response.HTTPResponse(response.ParamHTTPResponse{
				Code: http.StatusInternalServerError,
				Err:  err,
				Data: nil,
				Gin:  c,
			})
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "ユーザーを作成しました。",
		Gin:     c,
	})
}
