package user

import (
	userrepository "cleaning/internal/repository/user"
	"context"
)

type service struct {
	userRepository userrepository.Repository
}

type Service interface {
	GetAll(context.Context) ([]GetUserResult, error)
}

func NewService(userRepository userrepository.Repository) Service {
	return &service{
		userRepository: userRepository,
	}
}

func (s *service) GetAll(ctx context.Context) ([]GetUserResult, error) {
	users, err := s.userRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []GetUserResult
	for _, user := range users {
		result = append(result, GetUserResult{
			UserID:      user.ID,
			UUID:        user.UUID,
			UserName:    user.UserName,
			Email:       user.Email,
			IsActive:    user.IsActive,
			LastLoginAt: user.LastLoginAt,
			RoleName:    user.Role.Name,
		})
	}
	return result, nil
}
