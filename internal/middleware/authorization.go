package middleware

import (
	"cleaning/internal/common/response"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/services/permission"
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
			response.InternalServerError(c)
			c.Abort()
			return
		}

		if !hasPermission(permissions, requiredPermission) {
			response.Forbidden(c, errConstant.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}

func hasPermission(permissions []string, required string) bool {
	return slices.Contains(permissions, required)
}
