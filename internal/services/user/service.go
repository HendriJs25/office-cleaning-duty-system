package user

import (
	"cleaning/internal/common/email"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	employeerepository "cleaning/internal/repository/employee"
	rolerepository "cleaning/internal/repository/role"
	sessionrepository "cleaning/internal/repository/session"
	userrepository "cleaning/internal/repository/user"
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userRepository     userrepository.Repository
	roleRepository     rolerepository.Repository
	employeeRepository employeerepository.Repository
	sessionRepository  sessionrepository.Repository
}

type Service interface {
	GetAll(context.Context) ([]GetUserResult, error)
	GetUserDetail(context.Context, uuid.UUID) (*GetUserDetailResult, error)
	Create(context.Context, CreateUserInput) error
	DeactivateUser(context.Context, uuid.UUID) error
	ActivateUser(context.Context, uuid.UUID) error
	UpdateUserByAdmin(context.Context, uuid.UUID, uuid.UUID, UpdateUserByAdminInput) error
}

func NewService(userRepository userrepository.Repository, roleRepository rolerepository.Repository, employeeRepository employeerepository.Repository, sessionRepository sessionrepository.Repository) Service {
	return &service{
		userRepository:     userRepository,
		roleRepository:     roleRepository,
		employeeRepository: employeeRepository,
		sessionRepository:  sessionRepository,
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
			ID:          user.ID,
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

func (s *service) GetUserDetail(ctx context.Context, uuid uuid.UUID) (*GetUserDetailResult, error) {
	user, err := s.userRepository.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	var employee *EmployeeInfo
	if user.Employee != nil {
		employee = &EmployeeInfo{
			ID:                  user.Employee.ID,
			UUID:                user.Employee.UUID,
			OfficeName:          user.Employee.Office.Name,
			FullName:            fmt.Sprintf("%s %s", user.Employee.FamilyName, user.Employee.GivenName),
			EmploymentStartDate: user.Employee.EmploymentStartDate,
		}
	}

	return &GetUserDetailResult{
		User: GetUserResult{
			ID:          user.ID,
			UUID:        user.UUID,
			UserName:    user.UserName,
			Email:       user.Email,
			IsActive:    user.IsActive,
			LastLoginAt: user.LastLoginAt,
			RoleName:    user.Role.Name,
		},
		Employee: employee,
	}, nil
}

func (s *service) Create(ctx context.Context, createInput CreateUserInput) error {
	role, err := s.roleRepository.FindByID(ctx, createInput.RoleID)
	if err != nil {
		return err
	}

	if !role.IsActive {
		return errConstant.ErrRoleInActive
	}

	if createInput.EmployeeID != nil {
		employee, err := s.employeeRepository.FindByID(ctx, createInput.EmployeeID)
		if err != nil {
			return err
		}

		if !employee.IsActive {
			return errConstant.ErrEmployeeInActive
		}
	}

	normalizedEmail := email.Normalize(createInput.Email)

	exist, err := s.userRepository.ExistByEmail(ctx, normalizedEmail)
	if err != nil {
		return err
	}
	if exist {
		return errConstant.ErrAlreadyExists
	}

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

func (s *service) DeactivateUser(ctx context.Context, uuid uuid.UUID) error {
	user, err := s.userRepository.FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	if !user.IsActive {
		return errConstant.ErrAlreadyDeactivated
	}

	if err := s.sessionRepository.DeleteByUserUUID(ctx, user.UUID); err != nil {
		return err
	}

	if err := s.userRepository.UpdateStatus(ctx, user.UUID, false); err != nil {
		return err
	}

	return nil
}

func (s *service) ActivateUser(ctx context.Context, uuid uuid.UUID) error {
	user, err := s.userRepository.FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	if user.IsActive {
		return errConstant.ErrAlreadyActivated
	}

	if err := s.userRepository.UpdateStatus(ctx, user.UUID, true); err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateUserByAdmin(ctx context.Context, actorUUID uuid.UUID, targetUUID uuid.UUID, update UpdateUserByAdminInput) error {
	if actorUUID == targetUUID {
		return errConstant.ErrCannotUpdateSelf
	}

	role, err := s.roleRepository.FindByID(ctx, update.RoleID)
	if err != nil {
		return err
	}
	if !role.IsActive {
		return errConstant.ErrRoleInActive
	}

	if update.EmployeeID != nil {
		employee, err := s.employeeRepository.FindByID(ctx, update.EmployeeID)
		if err != nil {
			return err
		}

		if !employee.IsActive {
			return errConstant.ErrEmployeeInActive
		}
	}

	user, err := s.userRepository.FindByUUID(ctx, targetUUID)
	if err != nil {
		return err
	}

	fields := map[string]any{
		"role_id":     update.RoleID,
		"employee_id": update.EmployeeID,
		"is_active":   update.IsActive,
	}

	err = s.userRepository.UpdateByUUID(ctx, targetUUID, fields)
	if err != nil {
		return err
	}

	if user.RoleID != update.RoleID || !update.IsActive {
		err := s.sessionRepository.DeleteByUserUUID(ctx, targetUUID)
		if err != nil {
			return err
		}
	}

	return nil
}
