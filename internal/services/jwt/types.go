package jwt

import (
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessToken struct {
	Value     string
	ExpiresAt time.Time
}

type claims struct {
	UserUUID uuid.UUID

	jwtlib.RegisteredClaims
}
