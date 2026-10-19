package office

import (
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	officerepository "cleaning/internal/repository/office"
	"context"
)

type service struct {
	officeRepository officerepository.Repository
}

type Service interface {
	Create(context.Context, CreateOfficeInput) error
	GetAllOffices(context.Context) ([]GetOfficeResult, error)
	GetAllActiveOffices(context.Context) ([]GetActiveOfficeResult, error)
}

func NewService(officeRepository officerepository.Repository) Service {
	return &service{
		officeRepository: officeRepository,
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
