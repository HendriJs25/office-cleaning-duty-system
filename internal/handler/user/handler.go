package user

import (
	"cleaning/internal/common/response"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/dto/request"
	responsedto "cleaning/internal/domain/dto/response"
	"cleaning/internal/logger"
	"cleaning/internal/middleware"
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

	var users []responsedto.UserInfoResponse
	for _, user := range result {
		users = append(users, responsedto.UserInfoResponse{
			ID:          user.ID,
			UUID:        user.UUID,
			UserName:    user.UserName,
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
		User: responsedto.UserInfoResponse{
			ID:          userResult.User.ID,
			UUID:        userResult.User.UUID,
			UserName:    userResult.User.UserName,
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

func (h *Handler) GetProfile(c *gin.Context) {
	identity, ok := middleware.RequireIdentity(c)
	if !ok {
		return
	}

	userProfileResult, err := h.UserService.GetProfile(c.Request.Context(), identity.UUID)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrNotFound):
			response.NotFound(c)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("user_uuid", identity.UUID).WithError(err).Error("failed to get user profile")
			response.InternalServerError(c)
			return
		}
	}

	var employeeProfile *responsedto.EmployeeProfileResponse
	if userProfileResult.Employee != nil {
		employeeProfile = &responsedto.EmployeeProfileResponse{
			OfficeName:          userProfileResult.Employee.OfficeName,
			FullName:            userProfileResult.Employee.FullName,
			FullNameKana:        userProfileResult.Employee.FullNameKana,
			EmploymentStartDate: userProfileResult.Employee.EmploymentStartDate,
		}
	}

	userProfile := responsedto.UserProfileResponse{
		UserName: userProfileResult.UserName,
		Email:    userProfileResult.Email,
		RoleName: userProfileResult.RoleName,
		Employee: employeeProfile,
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: userProfile,
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
		case errors.Is(err, errConstant.ErrEmployeeAlreadyAssigned):
			response.Conflict(c, err)
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
		Code:    http.StatusCreated,
		Message: "ユーザーを作成しました",
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

func (h *Handler) UpdateUserByAdmin(c *gin.Context) {
	identity, ok := middleware.RequireIdentity(c)
	if !ok {
		return
	}

	uuidStr := c.Param("uuid")
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		response.BadRequest(c)
		return
	}

	var req request.UpdateUserByAdminRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	err = h.UserService.UpdateUserByAdmin(c.Request.Context(), identity.UUID, parsedUUID, userservice.UpdateUserByAdminInput{
		RoleID:     req.RoleID,
		EmployeeID: req.EmployeeID,
		IsActive:   *req.IsActive,
	})

	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrCannotUpdateSelf):
			response.Forbidden(c, err)
			return
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
		case errors.Is(err, errConstant.ErrEmployeeAlreadyAssigned):
			response.Conflict(c, err)
			return
		case errors.Is(err, errConstant.ErrNotFound):
			response.NotFound(c)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("user_uuid", parsedUUID).WithError(err).Error("failed to update user")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "ユーザー情報を更新しました",
		Gin:     c,
	})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	identity, ok := middleware.RequireIdentity(c)
	if !ok {
		return
	}

	var req request.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	err := h.UserService.UpdateProfile(c.Request.Context(), identity.UUID, userservice.UpdateProfileInput{
		UserName: req.UserName,
		Email:    req.Email,
	})

	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrNotFound):
			response.NotFound(c)
			return
		case errors.Is(err, errConstant.ErrAlreadyExists):
			response.ConflictDuplicate(c, err, "メールアドレス")
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("user_uuid", identity.UUID).WithError(err).Error("failed to update profile")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "プロフィールを更新しました",
		Gin:     c,
	})
}

func (h *Handler) UpdatePassword(c *gin.Context) {
	identity, ok := middleware.RequireIdentity(c)
	if !ok {
		return
	}

	var req request.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	err := h.UserService.UpdatePassword(c.Request.Context(), identity.UUID, userservice.UpdatePasswordInput{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})

	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrNotFound):
			response.NotFound(c)
			return
		case errors.Is(err, errConstant.ErrPasswordIncorrect):
			response.UnprocessableEntity(c, err)
			return
		case errors.Is(err, errConstant.ErrPasswordIsSame):
			response.UnprocessableEntity(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("user_uuid", identity.UUID).WithError(err).Error("failed to update password")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "パスワードを更新しました",
		Gin:     c,
	})
}
