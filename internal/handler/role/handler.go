package role

import (
	"cleaning/internal/common/response"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/dto/request"
	responsedto "cleaning/internal/domain/dto/response"
	"cleaning/internal/logger"
	roleservice "cleaning/internal/services/role"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	customValidator "github.com/go-playground/validator/v10"
)

type Handler struct {
	roleService roleservice.Service
	validate    *customValidator.Validate
}

func NewHandler(roleService roleservice.Service, validate *customValidator.Validate) *Handler {
	return &Handler{
		roleService: roleService,
		validate:    validate,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	err := h.roleService.Create(c.Request.Context(), roleservice.CreateRoleInput{
		Code:            req.Code,
		Name:            req.Name,
		PermissionsCode: req.PermissionCodes,
	})

	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrAlreadyExists):
			response.ConflictDuplicate(c, err, "コード")
			return
		case errors.Is(err, errConstant.ErrInvalidPermission):
			response.UnprocessableEntity(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).Error("create role failed")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusCreated,
		Message: "ロールを作成しました",
		Gin:     c,
	})

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
