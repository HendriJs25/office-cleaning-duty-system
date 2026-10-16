package permission

import (
	"cleaning/internal/common/response"
	responsedto "cleaning/internal/domain/dto/response"
	"cleaning/internal/logger"
	permissionservice "cleaning/internal/services/permission"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	permissionService permissionservice.Service
}

func NewHandler(permissionService permissionservice.Service) *Handler {
	return &Handler{
		permissionService: permissionService,
	}
}

func (h *Handler) GetAllPermissionsForOptions(c *gin.Context) {
	permissions, err := h.permissionService.GetAllPermissionsForOptions(c.Request.Context())
	if err != nil {
		logger.WithContext(c.Request.Context()).WithError(err).Error("get all permissions for options failed")
		response.InternalServerError(c)
		return
	}

	var result []responsedto.GetAllPermissionsForOptionsResponse
	for _, permission := range permissions {
		result = append(result, responsedto.GetAllPermissionsForOptionsResponse{
			Code: permission.Code,
			Name: permission.Name,
		})
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}
