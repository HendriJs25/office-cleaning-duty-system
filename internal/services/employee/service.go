package employee

import (
	"cleaning/internal/common/name"
	"cleaning/internal/constants"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	employeerepository "cleaning/internal/repository/employee"
	officerepository "cleaning/internal/repository/office"
	verificationtokenrepository "cleaning/internal/repository/verificationtoken"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type service struct {
	employeeRepository          employeerepository.Repository
	officeRepository            officerepository.Repository
	verificationTokenRepository verificationtokenrepository.Repository
}

type Service interface {
	Create(context.Context, CreateEmployeeInput) error
	GetAllEmployees(context.Context) ([]GetEmployeeResult, error)
	GetAllAssignableEmployees(context.Context, *uuid.UUID) ([]GetAssignableEmployeesResult, error)
	GetEmployeeDetail(context.Context, uuid.UUID) (*GetEmployeeDetailsResult, error)
	UpdateEmployeeByAdmin(context.Context, uuid.UUID, UpdateEmployeeByAdminInput) error
}

func NewService(employeeRepository employeerepository.Repository, officeRepository officerepository.Repository, verificationTokenRepository verificationtokenrepository.Repository) Service {
	return &service{
		employeeRepository:          employeeRepository,
		officeRepository:            officeRepository,
		verificationTokenRepository: verificationTokenRepository,
	}
}

func (s *service) Create(ctx context.Context, input CreateEmployeeInput) error {
	office, err := s.officeRepository.GetByID(ctx, input.OfficeID)
	if err != nil {
		return err
	}

	if !office.IsActive {
		return errConstant.ErrInActive
	}

	employee := &model.Employee{
		UUID:                uuid.New(),
		OfficeID:            input.OfficeID,
		FamilyName:          input.FamilyName,
		GivenName:           input.GivenName,
		FamilyNameKana:      input.FamilyNameKana,
		GivenNameKana:       input.GivenNameKana,
		IsActive:            true,
		EmploymentStartDate: input.EmploymentStartDate,
		EmploymentEndDate:   input.EmploymentEndDate,
	}

	if err := s.employeeRepository.Create(ctx, employee); err != nil {
		return err
	}

	return nil
}

func (s *service) GetAllEmployees(ctx context.Context) ([]GetEmployeeResult, error) {
	employees, err := s.employeeRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(employees) == 0 {
		return []GetEmployeeResult{}, nil
	}

	var result []GetEmployeeResult
	for _, employee := range employees {
		result = append(result, GetEmployeeResult{
			ID:                  employee.ID,
			UUID:                employee.UUID,
			OfficeName:          employee.Office.Name,
			FullName:            name.BuildFullName(employee.FamilyName, employee.GivenName),
			IsActive:            employee.IsActive,
			EmploymentStartDate: employee.EmploymentStartDate,
		})
	}

	return result, nil
}

func (s *service) GetAllAssignableEmployees(ctx context.Context, userUUID *uuid.UUID) ([]GetAssignableEmployeesResult, error) {
	employees, err := s.employeeRepository.FindAssignableEmployees(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	if len(employees) == 0 {
		return []GetAssignableEmployeesResult{}, nil
	}

	var result []GetAssignableEmployeesResult
	var fullName string
	for _, employee := range employees {
		fullName = fmt.Sprintf("%s %s", employee.FamilyName, employee.GivenName)
		result = append(result, GetAssignableEmployeesResult{
			ID:       employee.ID,
			FullName: fullName,
		})
	}
	return result, nil
}

func (s *service) GetEmployeeDetail(ctx context.Context, employeeUUID uuid.UUID) (*GetEmployeeDetailsResult, error) {
	employee, err := s.employeeRepository.FindByUUID(ctx, employeeUUID)
	if err != nil {
		return nil, err
	}

	var invitationToken *Token

	token, err := s.verificationTokenRepository.FindByEmployeeIDAndTokenType(ctx, employee.ID, constants.TokenTypeUserInvitation)
	if err != nil {
		if !errors.Is(err, errConstant.ErrNotFound) {
			return nil, err
		}
	} else {
		invitationToken = &Token{
			TokenType: token.TokenType,
			ExpiresAt: token.ExpiresAt,
		}
	}

	result := &GetEmployeeDetailsResult{
		ID:                  employee.ID,
		UUID:                employeeUUID,
		OfficeName:          employee.Office.Name,
		FullName:            name.BuildFullName(employee.FamilyName, employee.GivenName),
		FullNameKana:        name.BuildFullNameKana(employee.FamilyNameKana, employee.GivenNameKana),
		IsActive:            employee.IsActive,
		EmploymentStartDate: employee.EmploymentStartDate,
		EmploymentEndDate:   employee.EmploymentEndDate,
		Token:               invitationToken,
	}

	return result, nil
}

func (s *service) UpdateEmployeeByAdmin(ctx context.Context, employeeUUID uuid.UUID, input UpdateEmployeeByAdminInput) error {
	office, err := s.officeRepository.GetByID(ctx, input.OfficeID)
	if err != nil {
		return err
	}

	if !office.IsActive {
		return errConstant.ErrInActive
	}

	updates := map[string]any{
		"office_id":             input.OfficeID,
		"family_name":           input.FamilyName,
		"given_name":            input.GivenName,
		"family_name_kana":      input.FamilyNameKana,
		"given_name_kana":       input.GivenNameKana,
		"is_active":             input.IsActive,
		"employment_start_date": input.EmploymentStartDate,
		"employment_end_date":   input.EmploymentEndDate,
	}

	err = s.employeeRepository.UpdateByUUID(ctx, employeeUUID, updates)
	if err != nil {
		return err
	}
	return nil
}
