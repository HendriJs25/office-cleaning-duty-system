package user

import (
	"cleaning/internal/common/email"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	rolerepository "cleaning/internal/repository/role"
	userrepository "cleaning/internal/repository/user"
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userRepository userrepository.Repository
	roleRepository rolerepository.Repository
}

type Service interface {
	GetAll(context.Context) ([]GetUserResult, error)
	Create(context.Context, CreateUserInput) error
}

func NewService(userRepository userrepository.Repository, roleRepository rolerepository.Repository) Service {
	return &service{
		userRepository: userRepository,
		roleRepository: roleRepository,
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

func (s *service) Create(ctx context.Context, createInput CreateUserInput) error {
	exist, err := s.roleRepository.ExistByID(ctx, createInput.RoleID)
	if err != nil {
		return err
	}
	if !exist {
		return errConstant.ErrNotFound
	}

	normalizedEmail := email.Normalize(createInput.Email)

	exist, err = s.userRepository.ExistByEmail(ctx, normalizedEmail)
	if err != nil {
		return err
	}
	if exist {
		return errConstant.ErrAlreadyExists
	}

	// todo: employee is exist function

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(createInput.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user := model.User{
		UUID:         uuid.New(),
		RoleID:       createInput.RoleID,
		EmployeeID:   createInput.EmployeeID,
		UserName:     createInput.UserName,
		Email:        normalizedEmail,
		PasswordHash: string(hashPassword),
		IsActive:     true,
	}

	err = s.userRepository.Create(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}
