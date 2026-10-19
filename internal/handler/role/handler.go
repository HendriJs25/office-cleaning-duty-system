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
	"strconv"
	"strings"

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

func (h *Handler) GetRoleDetail(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c)
		return
	}

	role, err := h.roleService.GetRoleDetail(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrRoleNotFound):
			response.NotFound(c)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("roleID", id).WithError(err).Error("Get role failed")
			response.InternalServerError(c)
			return
		}
	}

	permissionResponse := make([]responsedto.PermissionsResponse, 0, len(role.Permissions))
	for _, permission := range role.Permissions {
		permissionResponse = append(permissionResponse, responsedto.PermissionsResponse{
			ID:   permission.ID,
			Code: permission.Code,
			Name: permission.Name,
		})
	}

	roleResponse := responsedto.GetRoleDetailResponse{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		IsActive:    role.IsActive,
		Permissions: permissionResponse,
	}
	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: roleResponse,
		Gin:  c,
	})

}

func (h *Handler) DeactivateRole(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c)
		return
	}

	err = h.roleService.DeactivateRole(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrRoleNotFound):
			response.NotFound(c)
			return
		case errors.Is(err, errConstant.ErrSystemRoleImmutable):
			response.Conflict(c, err)
			return
		case errors.Is(err, errConstant.ErrAlreadyDeactivated):
			response.Conflict(c, err)
			return
		case errors.Is(err, errConstant.ErrRoleInUse):
			response.Conflict(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("id", id).WithError(err).Error("Deactivate role failed")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "ロールを無効化しました",
		Gin:     c,
	})
}

func (h *Handler) ActivateRole(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c)
		return
	}

	err = h.roleService.ActivateRole(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrRoleNotFound):
			response.NotFound(c)
			return
		case errors.Is(err, errConstant.ErrAlreadyActivated):
			response.Conflict(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("id", id).WithError(err).Error("Activate role failed")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "ロールを有効化しました",
		Gin:     c,
	})
}

func (h *Handler) UpdateRole(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c)
		return
	}

	var req request.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	err = h.roleService.UpdateRoleByID(c.Request.Context(), id, roleservice.UpdateRoleInput{
		Name:            req.Name,
		Description:     req.Description,
		PermissionsCode: req.PermissionCodes,
	})
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrRoleNotFound):
			response.NotFound(c)
			return
		case errors.Is(err, errConstant.ErrSystemRoleImmutable):
			response.Conflict(c, err)
			return
		case errors.Is(err, errConstant.ErrInvalidPermission):
			response.UnprocessableEntity(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("id", id).WithError(err).Error("Update role failed")
			response.InternalServerError(c)
			return
		}
	}
	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "ロールを更新しました",
		Gin:     c,
	})

}
