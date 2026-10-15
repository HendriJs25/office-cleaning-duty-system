package verificationtoken

import (
	"cleaning/internal/common/token"
	"cleaning/internal/constants"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	"cleaning/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type service struct {
	transactionManager repository.TransactionManager
}

type Service interface {
	CreateRegistrationToken(context.Context, uuid.UUID) (*CreateRegistrationTokenResult, error)
}

func NewService(transactionManager repository.TransactionManager) Service {
	return &service{
		transactionManager: transactionManager,
	}
}

func (s *service) CreateRegistrationToken(ctx context.Context, employeeUUID uuid.UUID) (*CreateRegistrationTokenResult, error) {
	var result CreateRegistrationTokenResult

	err := s.transactionManager.WithinTransaction(ctx, func(repositories *repository.Registry) error {
		employee, err := repositories.Employee.FindByUUIDForUpdate(ctx, employeeUUID)
		if err != nil {
			return err
		}

		if !employee.IsActive {
			return errConstant.ErrEmployeeInActive
		}

		assignable, err := repositories.Employee.IsAssignable(ctx, employeeUUID)
		if err != nil {
			return err
		}

		if !assignable {
			return errConstant.ErrEmployeeAlreadyAssigned
		}

		err = repositories.VerificationToken.DeleteByEmployeeIDAndTokenType(ctx, employee.ID, constants.TokenTypeUserInvitation)
		if err != nil {
			return err
		}

		rawToken := token.Generate()
		hashToken := token.Hash(rawToken)
		expiresAT := time.Now().UTC().Add(time.Hour * 24)

		err = repositories.VerificationToken.Create(ctx, &model.VerificationToken{
			EmployeeID: employee.ID,
			HashToken:  hashToken,
			TokenType:  constants.TokenTypeUserInvitation,
			ExpiresAt:  expiresAT,
		})

		if err != nil {
			return err
		}

		result = CreateRegistrationTokenResult{
			RawToken:  rawToken,
			ExpiresAt: expiresAT,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}
