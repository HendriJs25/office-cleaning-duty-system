package user

import (
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	sessionrepository "cleaning/internal/repository/session"
	userrepository "cleaning/internal/repository/user"
	"cleaning/internal/services/jwt"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userRepository    userrepository.Repository
	sessionRepository sessionrepository.Repository
	jwtService        jwt.Service
}

type Service interface {
	Authenticate(context.Context, AuthenticateInput) (*AuthenticatedUser, error)
	Login(context.Context, LoginInput) (*LoginResult, error)
}

func NewService(userRepository userrepository.Repository, sessionRepository sessionrepository.Repository, jwtService jwt.Service) Service {
	return &service{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		jwtService:        jwtService,
	}
}

func (s *service) Authenticate(ctx context.Context, input AuthenticateInput) (*AuthenticatedUser, error) {
	email := normalizeEmail(input.Email)

	user, err := s.userRepository.FindByEmail(ctx, email)
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

	return &LoginResult{
		Token: accessToken.Value,
		User:  *authenticatedUser,
	}, nil

}
