package user

import (
	"cleaning/internal/common/response"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/dto/request"
	responsedto "cleaning/internal/domain/dto/response"
	"cleaning/internal/logger"
	userservice "cleaning/internal/services/user"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	customValidator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
		logger.WithContext(c.Request.Context()).WithError(err).Error("Get all user failed")
		response.InternalServerError(c)
		return
	}

	var users []responsedto.UserResponse
	for _, user := range result {
		users = append(users, responsedto.UserResponse{
			ID:          user.ID,
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

func (h *Handler) GetUserDetail(c *gin.Context) {
	uuidStr := c.Param("uuid")

	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		response.BadRequest(c)
		return
	}

	userResult, err := h.UserService.GetUserDetail(c.Request.Context(), parsedUUID)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrNotFound):
			response.NotFound(c)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("user_uuid", parsedUUID).WithError(err).Error("failed to get user detail")
			response.InternalServerError(c)
			return
		}
	}

	var employeeInfo *responsedto.EmployeeInfoResponse
	if userResult.Employee != nil {
		employeeInfo = &responsedto.EmployeeInfoResponse{
			ID:                  userResult.Employee.ID,
			UUID:                userResult.Employee.UUID,
			OfficeName:          userResult.Employee.OfficeName,
			FullName:            userResult.Employee.FullName,
			EmploymentStartDate: userResult.Employee.EmploymentStartDate,
		}
	}

	userDetail := responsedto.UserDetailResponse{
		User: responsedto.UserResponse{
			ID:          userResult.User.ID,
			UUID:        userResult.User.UUID,
			Username:    userResult.User.UserName,
			Email:       userResult.User.Email,
			IsActive:    &userResult.User.IsActive,
			LastLoginAt: userResult.User.LastLoginAt,
			RoleName:    userResult.User.RoleName,
		},
		Employee: employeeInfo,
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: userDetail,
		Gin:  c,
	})
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
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
		case errors.Is(err, errConstant.ErrRoleNotFound):
			response.InvalidOption(c, err, "ロール")
			return
		case errors.Is(err, errConstant.ErrEmployeeNotFound):
			response.InvalidOption(c, err, "従業員")
			return
		case errors.Is(err, errConstant.ErrRoleInActive):
			response.InactiveOption(c, err, "ロール")
			return
		case errors.Is(err, errConstant.ErrEmployeeInActive):
			response.InactiveOption(c, err, "従業員")
			return
		case errors.Is(err, errConstant.ErrAlreadyExists):
			response.ConflictDuplicate(c, err, "メールアドレス")
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).Error("failed to create user")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "ユーザーを作成しました。",
		Gin:     c,
	})
}

func (h *Handler) DeactivateUser(c *gin.Context) {
	uuidStr := c.Param("uuid")
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		response.BadRequest(c)
		return
	}

	err = h.UserService.DeactivateUser(c.Request.Context(), parsedUUID)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrNotFound):
			response.NotFound(c)
			return
		case errors.Is(err, errConstant.ErrAlreadyDeactivated):
			response.Conflict(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("user_uuid", parsedUUID).WithError(err).Error("failed to deactivate user")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "ユーザーを無効化しました",
		Gin:     c,
	})
}

func (h *Handler) ActivateUser(c *gin.Context) {
	uuidStr := c.Param("uuid")
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		response.BadRequest(c)
		return
	}

	err = h.UserService.ActivateUser(c.Request.Context(), parsedUUID)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrNotFound):
			response.NotFound(c)
			return
		case errors.Is(err, errConstant.ErrAlreadyActivated):
			response.Conflict(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("user_uuid", parsedUUID).WithError(err).Error("failed to activate user")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "ユーザーを有効化しました",
		Gin:     c,
	})
}
