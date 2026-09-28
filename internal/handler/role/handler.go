package role

import (
	"cleaning/internal/common/response"
	errConstant "cleaning/internal/constants/error"
	responsedto "cleaning/internal/domain/dto/response"
	roleservice "cleaning/internal/services/role"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	roleService roleservice.Service
}

func NewHandler(roleService roleservice.Service) *Handler {
	return &Handler{
		roleService: roleService,
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	result, err := h.roleService.GetAll(c.Request.Context())
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrNotFound):
			response.HTTPResponse(response.ParamHTTPResponse{
				Code: http.StatusOK,
				Data: nil,
				Gin:  c,
			})
			return
		default:
			slog.Error("get all roles failed", "error", err)
			response.HTTPResponse(response.ParamHTTPResponse{
				Code: http.StatusInternalServerError,
				Err:  err,
				Data: nil,
				Gin:  c,
			})
			return
		}
	}

	var roles []responsedto.RoleResponse
	for _, role := range result {
		roles = append(roles, responsedto.RoleResponse{
			ID:       role.ID,
			Code:     role.Code,
			Name:     role.Name,
			IsActive: role.IsActive,
		})
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: roles,
		Gin:  c,
	})
}
