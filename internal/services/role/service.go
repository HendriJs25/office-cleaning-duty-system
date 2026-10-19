package role

import (
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	"cleaning/internal/repository"
	permissionrepository "cleaning/internal/repository/permission"
	"cleaning/internal/repository/permissioncache"
	rolerepository "cleaning/internal/repository/role"
	"context"
)

type service struct {
	roleRepository            rolerepository.Repository
	permissionRepository      permissionrepository.Repository
	permissionCacheRepository permissioncache.Repository
	transactionManager        repository.TransactionManager
}

type Service interface {
	Create(context.Context, CreateRoleInput) error
	GetAll(context.Context) ([]RoleResult, error)
	GetAllActiveRoles(context.Context) ([]GetAllActiveRolesResult, error)
	GetRoleDetail(context.Context, int64) (*GetRoleDetailResult, error)
	DeactivateRole(context.Context, int64) error
	ActivateRole(context.Context, int64) error
	UpdateRoleByID(context.Context, int64, UpdateRoleInput) error
}

func NewService(roleRepository rolerepository.Repository, permissionRepository permissionrepository.Repository, permissionCacheRepository permissioncache.Repository, transactionManager repository.TransactionManager) Service {
	return &service{
		roleRepository:            roleRepository,
		permissionRepository:      permissionRepository,
		permissionCacheRepository: permissionCacheRepository,
		transactionManager:        transactionManager,
	}
}

func (s *service) Create(ctx context.Context, input CreateRoleInput) error {
	exist, err := s.roleRepository.ExistByCode(ctx, input.Code)
	if err != nil {
		return err
	}
	if exist {
		return errConstant.ErrAlreadyExists
	}

	permissions, err := s.permissionRepository.FindByCodes(ctx, input.PermissionsCode)
	if err != nil {
		return err
	}

	if len(input.PermissionsCode) != len(permissions) {
		return errConstant.ErrInvalidPermission
	}

	err = s.transactionManager.WithinTransaction(ctx, func(repositories *repository.Registry) error {
		role := &model.Role{
			Code:     input.Code,
			Name:     input.Name,
			IsActive: true,
		}

		err := repositories.Role.Create(ctx, role)
		if err != nil {
			return err
		}

		err = repositories.Role.AssignPermission(ctx, role.ID, permissions)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAll(ctx context.Context) ([]RoleResult, error) {
	roles, err := s.roleRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return []RoleResult{}, nil
	}

	var result []RoleResult

	for _, role := range roles {
		result = append(result, RoleResult{
			ID:       role.ID,
			Code:     role.Code,
			Name:     role.Name,
			IsActive: role.IsActive,
		})
	}
	return result, nil
}

func (s *service) GetAllActiveRoles(ctx context.Context) ([]GetAllActiveRolesResult, error) {
	result, err := s.roleRepository.FindAllActive(ctx)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []GetAllActiveRolesResult{}, nil
	}

	var roles []GetAllActiveRolesResult
	for _, role := range result {
		roles = append(roles, GetAllActiveRolesResult{
			ID:   role.ID,
			Name: role.Name,
		})
	}

	return roles, nil
}

func (s *service) GetRoleDetail(ctx context.Context, roleID int64) (*GetRoleDetailResult, error) {
	role, err := s.roleRepository.FindByID(ctx, roleID, true)
	if err != nil {
		return nil, err
	}

	permissionResult := make([]PermissionInfo, 0, len(role.Permissions))
	for _, permission := range role.Permissions {
		permissionResult = append(permissionResult, PermissionInfo{
			ID:   permission.ID,
			Code: permission.Code,
			Name: permission.Name,
		})
	}

	roleResult := GetRoleDetailResult{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		IsActive:    role.IsActive,
		Permissions: permissionResult,
	}
	return &roleResult, nil
}

func (s *service) DeactivateRole(ctx context.Context, roleID int64) error {
	role, err := s.roleRepository.FindByID(ctx, roleID, false)
	if err != nil {
		return err
	}

	if role.IsSystem {
		return errConstant.ErrSystemRoleImmutable
	}

	if !role.IsActive {
		return errConstant.ErrAlreadyDeactivated
	}

	hasAssignedUser, err := s.roleRepository.HasAssignedUsers(ctx, roleID)
	if err != nil {
		return err
	}

	if hasAssignedUser {
		return errConstant.ErrRoleInUse
	}

	err = s.roleRepository.UpdateStatus(ctx, roleID, false)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ActivateRole(ctx context.Context, roleID int64) error {
	role, err := s.roleRepository.FindByID(ctx, roleID, false)
	if err != nil {
		return err
	}
	if role.IsActive {
		return errConstant.ErrAlreadyActivated
	}

	err = s.roleRepository.UpdateStatus(ctx, roleID, true)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateRoleByID(ctx context.Context, roleID int64, update UpdateRoleInput) error {
	role, err := s.roleRepository.FindByID(ctx, roleID, false)
	if err != nil {
		return err
	}
	if role.IsSystem {
		return errConstant.ErrSystemRoleImmutable
	}

	err = s.permissionCacheRepository.DeleteRolePermissions(ctx, roleID)
	if err != nil {
		return err
	}

	err = s.transactionManager.WithinTransaction(ctx, func(repositories *repository.Registry) error {
		role, err := repositories.Role.FindByIDForUpdate(ctx, roleID)
		if err != nil {
			return err
		}
		permissions, err := repositories.Permission.FindByCodes(ctx, update.PermissionsCode)
		if err != nil {
			return err
		}

		if len(update.PermissionsCode) != len(permissions) {
			return errConstant.ErrInvalidPermission
		}

		updatesRole := map[string]any{
			"name":        update.Name,
			"description": update.Description,
		}
		err = repositories.Role.UpdateByID(ctx, role.ID, updatesRole)
		if err != nil {
			return err
		}

		err = repositories.Role.UpdatePermissionByID(ctx, role.ID, permissions)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
