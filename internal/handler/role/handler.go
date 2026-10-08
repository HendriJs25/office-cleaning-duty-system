package role

import (
	"cleaning/internal/common/response"
	responsedto "cleaning/internal/domain/dto/response"
	"cleaning/internal/logger"
	roleservice "cleaning/internal/services/role"
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
		logger.WithContext(c.Request.Context()).WithError(err).Error("Get all roles failed")
		response.InternalServerError(c)
		return
	}

	if len(result) == 0 {
		response.Empty(c, "ロール")
		return
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

func (h *Handler) GetAllActiveRoles(c *gin.Context) {
	result, err := h.roleService.GetAllActiveRoles(c.Request.Context())
	if err != nil {
		logger.WithContext(c.Request.Context()).WithError(err).Error("Get all active roles failed")
		response.InternalServerError(c)
		return
	}

	if len(result) == 0 {
		response.Empty(c, "ロール")
		return
	}

	var roles []responsedto.GetActiveRolesResponse
	for _, role := range result {
		roles = append(roles, responsedto.GetActiveRolesResponse{
			ID:   role.ID,
			Name: role.Name,
		})
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: roles,
		Gin:  c,
	})

}
