package user

import (
	"cleaning/internal/common/response"
	responsedto "cleaning/internal/domain/dto/response"
	userservice "cleaning/internal/services/user"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	UserService userservice.Service
}

func NewHandler(userService userservice.Service) *Handler {
	return &Handler{
		UserService: userService,
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
