package role

import (
	rolerepository "cleaning/internal/repository/role"
	"context"
)

type service struct {
	roleRepository rolerepository.Repository
}

type Service interface {
	GetAll(context.Context) ([]RoleResult, error)
}

func NewService(roleRepository rolerepository.Repository) Service {
	return &service{
		roleRepository: roleRepository,
	}
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
