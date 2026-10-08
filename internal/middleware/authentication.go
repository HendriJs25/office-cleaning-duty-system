package middleware

import (
	"cleaning/internal/common/response"
	"cleaning/internal/constants"
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/logger"
	sessionrepository "cleaning/internal/repository/session"
	"cleaning/internal/services/jwt"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const authenticatedIdentityKey = "authenticated_identity"

type Identity struct {
	UUID     uuid.UUID
	UserName string
	Email    string
	RoleID   int64
	Token    string
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
			response.Unauthorized(c, errConstant.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, err := a.jwtService.ValidateAccessToken(token)
		if err != nil {
			response.Unauthorized(c, errConstant.ErrUnauthorized)
			c.Abort()
			return
		}

		session, err := a.sessionRepository.Get(c.Request.Context(), token)
		if err != nil {
			if errors.Is(err, errConstant.ErrNotFound) {
				response.Unauthorized(c, errConstant.ErrUnauthorized)
				c.Abort()
				return
			}

			logger.WithContext(c.Request.Context()).WithError(err).Error("failed to get session")
			response.InternalServerError(c)
			c.Abort()
			return
		}

		if session.UUID != claims.UserUUID {
			logger.WithContext(c.Request.Context()).WithFields(logrus.Fields{
				"jwt_user_uuid":     claims.UserUUID,
				"session_user_uuid": session.UUID,
			}).Error("authentication identity mismatch")
			response.Unauthorized(c, errConstant.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set(authenticatedIdentityKey, Identity{
			UUID:     session.UUID,
			UserName: session.UserName,
			Email:    session.Email,
			RoleID:   session.RoleID,
			Token:    token,
		})

		c.Next()
	}
}

func RequireIdentity(c *gin.Context) (Identity, bool) {
	value, exists := c.Get(authenticatedIdentityKey)
	if !exists {
		logger.WithContext(c.Request.Context()).Error("authenticated identity missing from context")
		response.InternalServerError(c)
		c.Abort()
		return Identity{}, false
	}

	identity, ok := value.(Identity)
	if !ok {
		logger.WithContext(c.Request.Context()).Error("invalid authenticated identity type")
		response.InternalServerError(c)
		c.Abort()
		return Identity{}, false
	}

	return identity, true
}
