package middleware

import (
	"cleaning/internal/common/response"
	"cleaning/internal/constants"
	"cleaning/internal/services/permission"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

type Authorization struct {
	permissionService permission.Service
}

func NewAuthorization(permissionService permission.Service) *Authorization {
	return &Authorization{
		permissionService: permissionService,
	}
}

func (a *Authorization) RequirePermission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		identity, ok := RequireIdentity(c)
		if !ok {
			return
		}

		permissions, err := a.permissionService.GetCodesByRoleID(c.Request.Context(), identity.RoleID)
		if err != nil {
			abortInternalServerError(c)
			return
		}

		if !hasPermission(permissions, requiredPermission) {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Response{
				Status:  constants.Error,
				Message: "この操作を実行する権限がありません。",
			})
			return
		}
		c.Next()
	}
}

func hasPermission(permissions []string, required string) bool {
	return slices.Contains(permissions, required)
}
