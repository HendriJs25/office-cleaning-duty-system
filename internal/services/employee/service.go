package employee

import (
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	employeerepository "cleaning/internal/repository/employee"
	officerepository "cleaning/internal/repository/office"
	"context"

	"github.com/google/uuid"
)

type service struct {
	employeeRepository employeerepository.Repository
	officeRepository   officerepository.Repository
}

type Service interface {
	Create(context.Context, CreateEmployeeInput) error
}

func NewService(employeeRepository employeerepository.Repository, officeRepository officerepository.Repository) Service {
	return &service{
		employeeRepository: employeeRepository,
		officeRepository:   officeRepository,
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
