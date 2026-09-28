package middleware

import (
	"cleaning/internal/common/response"
	"cleaning/internal/constants"
	errConstant "cleaning/internal/constants/error"
	sessionrepository "cleaning/internal/repository/session"
	"cleaning/internal/services/jwt"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const authenticatedIdentityKey = "authenticated_identity"

type Identity struct {
	UUID     uuid.UUID
	UserName string
	Email    string
	RoleID   int64
}

type Authentication struct {
	jwtService        jwt.Service
	sessionRepository sessionrepository.Repository
}

func NewAuthentication(jwtService jwt.Service, sessionRepository sessionrepository.Repository) *Authentication {
	return &Authentication{
		jwtService:        jwtService,
		sessionRepository: sessionRepository,
	}
}

func (a *Authentication) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(constants.AccessTokenName)
		if err != nil {
			abortUnauthorized(c)
			return
		}

		claims, err := a.jwtService.ValidateAccessToken(token)
		if err != nil {
			abortUnauthorized(c)
			return
		}

		session, err := a.sessionRepository.Get(c.Request.Context(), token)
		if err != nil {
			if errors.Is(err, errConstant.ErrNotFound) {
				slog.Error("session not found")
			}
			abortUnauthorized(c)
			return
		}

		if session.UUID != claims.UserUUID {
			slog.Error("authentication identity mismatch", "jwt_user_uuid", claims.UserUUID, "session_user_uuid", session.UUID)
			abortUnauthorized(c)
			return
		}

		c.Set(authenticatedIdentityKey, Identity{
			UUID:     session.UUID,
			UserName: session.UserName,
			Email:    session.Email,
			RoleID:   session.RoleID,
		})

		c.Next()
	}
}

func IdentityFromContext(c *gin.Context) (Identity, bool) {
	value, exists := c.Get(authenticatedIdentityKey)
	if !exists {
		return Identity{}, false
	}

	identity, ok := value.(Identity)
	if !ok {
		return Identity{}, false
	}

	return identity, true
}

func abortUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
		Status:  constants.Error,
		Message: "認証が必要です。",
		Data:    nil,
	})
}
