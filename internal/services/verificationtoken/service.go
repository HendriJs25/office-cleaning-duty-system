package verificationtoken

import (
	"cleaning/internal/common/name"
	"cleaning/internal/common/token"
	"cleaning/internal/constants"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	"cleaning/internal/repository"
	employeerepository "cleaning/internal/repository/employee"
	verificationtokenrepository "cleaning/internal/repository/verificationtoken"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type service struct {
	transactionManager          repository.TransactionManager
	verificationTokenRepository verificationtokenrepository.Repository
	employeeRepository          employeerepository.Repository
}

type Service interface {
	CreateRegistrationToken(context.Context, uuid.UUID) (*CreateRegistrationTokenResult, error)
	VerifyRegistrationToken(context.Context, string) (*VerifyRegistrationTokenResult, error)
}

func NewService(transactionManager repository.TransactionManager, verificationTokenRepository verificationtokenrepository.Repository, employeeRepository employeerepository.Repository) Service {
	return &service{
		transactionManager:          transactionManager,
		verificationTokenRepository: verificationTokenRepository,
		employeeRepository:          employeeRepository,
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

func (s *service) VerifyRegistrationToken(ctx context.Context, rawToken string) (*VerifyRegistrationTokenResult, error) {
	tokenHash := token.Hash(rawToken)
	tokenResult, err := s.verificationTokenRepository.FindByTokenHashAndTokenType(ctx, tokenHash, constants.TokenTypeUserInvitation)
	if err != nil {
		if errors.Is(err, errConstant.ErrNotFound) {
			return nil, errConstant.ErrInvalidToken
		}
		return nil, err
	}

	if !time.Now().UTC().Before(tokenResult.ExpiresAt) {
		return nil, errConstant.ErrExpiredToken
	}

	employee, err := s.employeeRepository.FindByID(ctx, &tokenResult.EmployeeID)
	if err != nil {
		return nil, err
	}

	if !employee.IsActive {
		return nil, errConstant.ErrEmployeeInActive
	}

	assignable, err := s.employeeRepository.IsAssignable(ctx, employee.UUID)
	if err != nil {
		return nil, err
	}
	if !assignable {
		return nil, errConstant.ErrEmployeeAlreadyAssigned
	}

	result := &VerifyRegistrationTokenResult{
		EmployeeUUID: employee.UUID,
		FullName:     name.BuildFullName(employee.FamilyName, employee.GivenName),
		FullNameKana: name.BuildFullNameKana(employee.FamilyNameKana, employee.GivenNameKana),
		OfficeName:   employee.Office.Name,
	}

	return result, nil

}
