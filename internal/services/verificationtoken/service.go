package verificationtoken

import (
	employeerepository "cleaning/internal/repository/employee"
	verificationtokenrepository "cleaning/internal/repository/verificationtoken"
	"context"

	"github.com/google/uuid"
)

type service struct {
	verificationTokenRepository verificationtokenrepository.Repository
	employeeRepository          employeerepository.Repository
}

type Service interface {
	Create(context.Context, uuid.UUID) (*CreateRegistrationTokenResult, error)
}

func NewService(verificationTokenRepository verificationtokenrepository.Repository, employeeRepository employeerepository.Repository) Service {
	return &service{
		verificationTokenRepository: verificationTokenRepository,
		employeeRepository:          employeeRepository,
	}
}

func (s *service) Create(ctx context.Context, employeeUUID uuid.UUID) (*CreateRegistrationTokenResult, error) {
	employee, err := s.employeeRepository.FindByUUID(ctx, employeeUUID)
	if err != nil {
		return nil, err
	}

	if employee
}
