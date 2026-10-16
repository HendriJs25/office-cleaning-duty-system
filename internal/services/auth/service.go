package auth

import (
	"cleaning/internal/common/email"
	"cleaning/internal/common/token"
	"cleaning/internal/constants"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	"cleaning/internal/repository"
	rolerepository "cleaning/internal/repository/role"
	sessionrepository "cleaning/internal/repository/session"
	userrepository "cleaning/internal/repository/user"
	verificationtokenrepository "cleaning/internal/repository/verificationtoken"
	"cleaning/internal/services/jwt"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userRepository              userrepository.Repository
	sessionRepository           sessionrepository.Repository
	roleRepository              rolerepository.Repository
	verificationTokenRepository verificationtokenrepository.Repository
	jwtService                  jwt.Service
	transactionManager          repository.TransactionManager
}

type Service interface {
	Authenticate(context.Context, AuthenticateInput) (*AuthenticatedUser, error)
	Login(context.Context, LoginInput) (*LoginResult, error)
	Logout(context.Context, string) error
	Register(context.Context, RegisterInput) error
}

func NewService(userRepository userrepository.Repository,
	sessionRepository sessionrepository.Repository,
	roleRepository rolerepository.Repository,
	verificationTokenRepository verificationtokenrepository.Repository,
	jwtService jwt.Service,
	transactionManager repository.TransactionManager) Service {
	return &service{
		userRepository:              userRepository,
		sessionRepository:           sessionRepository,
		roleRepository:              roleRepository,
		verificationTokenRepository: verificationTokenRepository,
		jwtService:                  jwtService,
		transactionManager:          transactionManager,
	}
}

func (s *service) Authenticate(ctx context.Context, input AuthenticateInput) (*AuthenticatedUser, error) {
	normalizedEmail := email.Normalize(input.Email)

	user, err := s.userRepository.FindByEmail(ctx, normalizedEmail)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, errConstant.ErrPasswordIncorrect
	}

	if !user.IsActive {
		return nil, errConstant.ErrAccountIsDeactivated
	}

	return &AuthenticatedUser{
		ID:       user.ID,
		UUID:     user.UUID,
		Email:    user.Email,
		UserName: user.UserName,
		RoleID:   user.RoleID,
		RoleName: user.Role.Name,
	}, nil

}

func (s *service) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	authenticatedUser, err := s.Authenticate(ctx, AuthenticateInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtService.GenerateAccessToken(authenticatedUser.UUID)
	if err != nil {
		return nil, err
	}

	sessionCreatedAt := time.Now().UTC()

	sessionTTL := accessToken.ExpiresAt.Sub(sessionCreatedAt)

	session := model.Session{
		UUID:      authenticatedUser.UUID,
		UserName:  authenticatedUser.UserName,
		Email:     authenticatedUser.Email,
		RoleID:    authenticatedUser.RoleID,
		CreatedAt: sessionCreatedAt,
	}

	if err := s.sessionRepository.Set(ctx, accessToken.Value, session, sessionTTL); err != nil {
		return nil, err
	}

	err = s.userRepository.UpdateLastLoginAt(ctx, authenticatedUser.UUID)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken: accessToken,
		User:        *authenticatedUser,
	}, nil

}

func (s *service) Logout(ctx context.Context, accessToken string) error {
	return s.sessionRepository.Delete(ctx, accessToken)
}

func (s *service) Register(ctx context.Context, input RegisterInput) error {
	normalizedEmail := email.Normalize(input.Email)
	hashedToken := token.Hash(input.Token)

	tokenResult, err := s.verificationTokenRepository.FindByTokenHashAndTokenType(ctx, hashedToken, constants.TokenTypeUserInvitation)
	if err != nil {
		if errors.Is(err, errConstant.ErrNotFound) {
			return errConstant.ErrInvalidRegistrationToken
		}
		return err
	}

	exist, err := s.userRepository.ExistByEmail(ctx, normalizedEmail)
	if err != nil {
		return err
	}

	if exist {
		return errConstant.ErrAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	err = s.transactionManager.WithinTransaction(ctx, func(repositories *repository.Registry) error {
		employee, err := repositories.Employee.FindByIDForUpdate(ctx, tokenResult.EmployeeID)
		if err != nil {
			if errors.Is(err, errConstant.ErrNotFound) {
				return errConstant.ErrEmployeeNotFound
			}
			return err
		}

		if !employee.IsActive {
			return errConstant.ErrEmployeeInActive
		}

		assignable, err := repositories.Employee.IsAssignable(ctx, employee.UUID)
		if err != nil {
			return err
		}

		if !assignable {
			return errConstant.ErrEmployeeAlreadyAssigned
		}

		tokenResult, err = repositories.VerificationToken.FindByTokenHashAndTokenType(ctx, hashedToken, constants.TokenTypeUserInvitation)
		if err != nil {
			if errors.Is(err, errConstant.ErrNotFound) {
				return errConstant.ErrInvalidRegistrationToken
			}
			return err
		}

		if !time.Now().UTC().Before(tokenResult.ExpiresAt) {
			return errConstant.ErrExpiredToken
		}

		viewerRole, err := repositories.Role.FindByCode(ctx, constants.ViewerRoleCode)
		if err != nil {
			if errors.Is(err, errConstant.ErrNotFound) {
				return errConstant.ErrRoleNotFound
			}
			return err
		}

		if !viewerRole.IsActive {
			return errConstant.ErrRoleInActive
		}

		err = repositories.User.Create(ctx, &model.User{
			UUID:         uuid.New(),
			RoleID:       viewerRole.ID,
			EmployeeID:   &employee.ID,
			UserName:     input.UserName,
			Email:        normalizedEmail,
			PasswordHash: string(hashedPassword),
			IsActive:     true,
		})

		if err != nil {
			return err
		}

		err = repositories.VerificationToken.DeleteByEmployeeIDAndTokenType(ctx, tokenResult.EmployeeID, constants.TokenTypeUserInvitation)
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
