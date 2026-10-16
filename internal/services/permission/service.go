package permission

import (
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/logger"
	permissionrepository "cleaning/internal/repository/permission"
	permissioncacherepository "cleaning/internal/repository/permissioncache"
	"context"
	"errors"
)

type service struct {
	permissionRepository      permissionrepository.Repository
	permissionCacheRepository permissioncacherepository.Repository
}

type Service interface {
	GetAllPermissionsForOptions(context.Context) ([]GetPermissionsForOptionsResult, error)
	GetCodesByRoleID(context.Context, int64) ([]string, error)
}

func NewService(permissionRepository permissionrepository.Repository, permissionCacheRepository permissioncacherepository.Repository) Service {
	return &service{
		permissionRepository:      permissionRepository,
		permissionCacheRepository: permissionCacheRepository,
	}
}

func (s *service) GetAllPermissionsForOptions(ctx context.Context) ([]GetPermissionsForOptionsResult, error) {
	permissions, err := s.permissionRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var result []GetPermissionsForOptionsResult
	for _, permission := range permissions {
		result = append(result, GetPermissionsForOptionsResult{
			Code: permission.Code,
			Name: permission.Name,
		})
	}
	return result, nil
}

func (s *service) GetCodesByRoleID(ctx context.Context, roleID int64) ([]string, error) {
	codes, err := s.permissionCacheRepository.GetRolePermissions(ctx, roleID)
	if err == nil {
		return codes, nil
	}

	if !errors.Is(err, errConstant.ErrNotFound) {
		logger.WithContext(ctx).WithField("role_id", roleID).WithError(err).Error("failed to get cached permission codes by role ID")
	}

	codes, err = s.permissionRepository.FindCodesByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	if err := s.permissionCacheRepository.SetRolePermissions(ctx, roleID, codes); err != nil {
		logger.WithContext(ctx).WithField("role_id", roleID).WithError(err).Error("failed to set cached permission codes by role ID")
	}

	return codes, nil
}
