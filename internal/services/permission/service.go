package permission

import (
	errConstant "cleaning/internal/constants/error"
	permissionrepository "cleaning/internal/repository/permission"
	permissioncacherepository "cleaning/internal/repository/permissioncache"
	"context"
	"errors"
	"log/slog"
)

type service struct {
	permissionRepository      permissionrepository.Repository
	permissionCacheRepository permissioncacherepository.Repository
}

type Service interface {
	GetCodesByRoleID(context.Context, int64) ([]string, error)
}

func NewService(permissionRepository permissionrepository.Repository, permissionCacheRepository permissioncacherepository.Repository) Service {
	return &service{
		permissionRepository:      permissionRepository,
		permissionCacheRepository: permissionCacheRepository,
	}
}

func (s *service) GetCodesByRoleID(ctx context.Context, roleID int64) ([]string, error) {
	codes, err := s.permissionCacheRepository.GetRolePermissions(ctx, roleID)
	if err == nil {
		return codes, nil
	}

	if !errors.Is(err, errConstant.ErrNotFound) {
		slog.Error("get cache code permissions by role id failed ", "err", err)
	}

	codes, err = s.permissionRepository.FindCodesByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	if err := s.permissionCacheRepository.SetRolePermissions(ctx, roleID, codes); err != nil {
		slog.Error("set cache role permissions", "error", err)
	}

	return codes, nil
}
