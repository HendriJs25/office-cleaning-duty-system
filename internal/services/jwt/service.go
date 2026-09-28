package jwt

import (
	"cleaning/internal/config"
	errConstant "cleaning/internal/constants/error"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type service struct {
	secretKey      []byte
	issuer         string
	accessTokenTTL time.Duration
}

type Service interface {
	GenerateAccessToken(uuid.UUID) (*AccessToken, error)
	ValidateAccessToken(string) (*AccessTokenClaims, error)
}

func NewService(cfg *config.JWT) Service {
	return &service{
		secretKey:      []byte(cfg.SecretKey),
		issuer:         cfg.Issuer,
		accessTokenTTL: cfg.AccessTokenTTl,
	}
}

func (s *service) GenerateAccessToken(userUUID uuid.UUID) (*AccessToken, error) {
	now := time.Now()

	expiresAt := now.Add(s.accessTokenTTL)

	tokenClaims := claims{
		UserUUID: userUUID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer:    s.issuer,
			IssuedAt:  jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(expiresAt),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, tokenClaims)

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	return &AccessToken{
		Value:     tokenString,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *service) ValidateAccessToken(accessToken string) (*AccessTokenClaims, error) {
	tokenClaims := &claims{}

	token, err := jwtlib.ParseWithClaims(
		accessToken,
		tokenClaims,
		func(token *jwtlib.Token) (any, error) { return s.secretKey, nil },
		jwtlib.WithValidMethods([]string{jwtlib.SigningMethodHS256.Alg()}),
		jwtlib.WithIssuer(s.issuer),
		jwtlib.WithExpirationRequired(),
		jwtlib.WithIssuedAt(),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: validate access token: %v", errConstant.ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("%w: access token is invalid", errConstant.ErrInvalidToken)
	}

	return &AccessTokenClaims{
		UserUUID:  tokenClaims.UserUUID,
		Issuer:    tokenClaims.Issuer,
		IssuedAt:  tokenClaims.IssuedAt.Time,
		ExpiresAt: tokenClaims.ExpiresAt.Time,
	}, nil
}
