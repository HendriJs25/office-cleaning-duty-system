package config

import (
	"fmt"
	"time"
)

const (
	defaultJWTAccessTokenTTL = 24 * time.Hour
	minimumJWTSecretBytes    = 32
)

type JWT struct {
	SecretKey      string
	Issuer         string
	AccessTokenTTl time.Duration
}

func loadJWT() (*JWT, error) {
	secretKey, err := getEnvRequired("JWT_SECRET_KEY")
	if err != nil {
		return nil, err
	}

	issuer, err := getEnvRequired("JWT_ISSUER")
	if err != nil {
		return nil, err
	}

	accessTokenTTL, err := getEnvDuration("JWT_EXPIRATION_TIME", defaultJWTAccessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &JWT{
		SecretKey:      secretKey,
		Issuer:         issuer,
		AccessTokenTTl: accessTokenTTL,
	}, nil
}

func (j *JWT) Validate() error {
	if len(j.SecretKey) < minimumJWTSecretBytes {
		return fmt.Errorf("JWT_SECRET_KEY must be at least %d bytes", minimumJWTSecretBytes)
	}

	if j.AccessTokenTTl <= 0 {
		return fmt.Errorf("JWT_EXPIRATION_TIME must be greater than zero")
	}

	return nil
}
