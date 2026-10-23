package office

import (
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	"cleaning/internal/repository"
	officerepository "cleaning/internal/repository/office"
	"context"
)

type service struct {
	officeRepository   officerepository.Repository
	transactionManager repository.TransactionManager
}

type Service interface {
	Create(context.Context, CreateOfficeInput) error
	GetAllOffices(context.Context) ([]GetOfficeResult, error)
	GetAllActiveOffices(context.Context) ([]GetActiveOfficeResult, error)
	UpdateByID(context.Context, int64, UpdateOfficeInput) error
	DeactivateOffice(context.Context, int64) error
}

func NewService(officeRepository officerepository.Repository, transactionManager repository.TransactionManager) Service {
	return &service{
		officeRepository:   officeRepository,
		transactionManager: transactionManager,
	}
}

func (s *service) Create(ctx context.Context, input CreateOfficeInput) error {
	exist, err := s.officeRepository.ExistByCode(ctx, input.Code)
	if err != nil {
		return err
	}
	if exist {
		return errConstant.ErrAlreadyExists
	}

	office := &model.Office{
		Code:     input.Code,
		Name:     input.Name,
		Address:  input.Address,
		IsActive: true,
	}

	err = s.officeRepository.Create(ctx, office)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetAllOffices(ctx context.Context) ([]GetOfficeResult, error) {
	offices, err := s.officeRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]GetOfficeResult, 0, len(offices))

	for _, office := range offices {
		result = append(result, GetOfficeResult{
			ID:       office.ID,
			Code:     office.Code,
			Name:     office.Name,
			Address:  office.Address,
			IsActive: office.IsActive,
		})
	}
	return result, nil
}

func (s *service) GetAllActiveOffices(ctx context.Context) ([]GetActiveOfficeResult, error) {
	offices, err := s.officeRepository.GetAllActive(ctx)
	if err != nil {
		return nil, err
	}
	if offices == nil {
		return []GetActiveOfficeResult{}, nil
	}

	var result []GetActiveOfficeResult
	for _, office := range offices {
		result = append(result, GetActiveOfficeResult{
			ID:   office.ID,
			Name: office.Name,
		})
	}

	return result, nil
}

func (s *service) UpdateByID(ctx context.Context, officeID int64, input UpdateOfficeInput) error {
	updates := map[string]any{
		"name":    input.Name,
		"address": input.Address,
	}

	err := s.officeRepository.UpdateByID(ctx, officeID, updates)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeactivateOffice(ctx context.Context, officeID int64) error {
	err := s.transactionManager.WithinTransaction(ctx, func(repositories *repository.Registry) error {
		office, err := repositories.Office.GetByIDForUpdate(ctx, officeID)
		if err != nil {
			return err
		}

		if !office.IsActive {
			return errConstant.ErrAlreadyDeactivated
		}

		hasAssignedEmployee, err := repositories.Office.HasAssignedEmployee(ctx, officeID)
		if err != nil {
			return err
		}
		if hasAssignedEmployee {
			return errConstant.ErrOfficeInUse
		}

		err = repositories.OfficeArea.DeactivateByOfficeID(ctx, officeID)
		if err != nil {
			return err
		}

		err = repositories.Office.UpdateStatus(ctx, office.ID, false)
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
