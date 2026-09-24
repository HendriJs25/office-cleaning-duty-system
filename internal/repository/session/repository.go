package session

import (
	"cleaning/internal/domain/model"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	redislib "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const sessionKeyPrefix = "session:"

type repository struct {
	client redislib.Cmdable
}

type Repository interface {
	Set(context.Context, string, model.Session, time.Duration) error
}

func NewRepository(client redislib.Cmdable) Repository {
	return &repository{
		client: client,
	}
}

func (r *repository) Set(ctx context.Context, accessToken string, session model.Session, ttl time.Duration) error {
	return nil
}

func sessionKey(accessToken string) string {
	digest := sha256.Sum256([]byte(accessToken))

	return sessionKeyPrefix + hex.EncodeToString(digest[:])
}

func userSessionKey(uuid uuid.UUID) string {
	return fmt.Sprintf("user_session:%d", uuid)
}
