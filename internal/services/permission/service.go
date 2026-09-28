package permission

import (
	permissionrepository "cleaning/internal/repository/permission"
	"context"
)

type service struct {
	permissionRepository permissionrepository.Repository
}

type Service interface {
	GetCodesByRoleID(context.Context, int64) ([]string, error)
}

func NewService(permissionRepository permissionrepository.Repository) Service {
	return &service{
		permissionRepository: permissionRepository,
	}
}

func (s *service) GetCodesByRoleID(ctx context.Context, roleID int64) ([]string, error) {
	codes, err := s.permissionRepository.FindCodesByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	return codes, nil
}
