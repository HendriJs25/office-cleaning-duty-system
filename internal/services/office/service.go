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
