package user

import (
	errConstant "cleaning/internal/constants/error"
	userrepository "cleaning/internal/repository/user"
	"cleaning/internal/services/jwt"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userRepository userrepository.Repository
	jwtService     jwt.Service
}

type Service interface {
	Authenticate(context.Context, AuthenticateInput) (*AuthenticatedUser, error)
	Login(context.Context, LoginInput) (*LoginResult, error)
}

func NewService(userRepository userrepository.Repository, jwtService jwt.Service) Service {
	return &service{
		userRepository: userRepository,
		jwtService:     jwtService,
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
		RoleCode: user.Role.Code,
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

	return &LoginResult{
		Token: accessToken.Value,
		User:  *authenticatedUser,
	}, nil

}
