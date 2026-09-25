package session

import (
	"cleaning/internal/domain/model"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	payload, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("encode session: %w", err)
	}

	sessionKey := sessionKey(accessToken)
	indexKey := userSessionKey(session.UUID)

	_, err = r.client.TxPipelined(ctx, func(pipe redislib.Pipeliner) error {
		pipe.Set(ctx, sessionKey, payload, ttl)
		pipe.SAdd(ctx, indexKey, sessionKey)
		pipe.Expire(ctx, indexKey, ttl)
		return nil
	})

	if err != nil {
		return fmt.Errorf("store session: %w", err)
	}

	return nil
}

func sessionKey(accessToken string) string {
	digest := sha256.Sum256([]byte(accessToken))

	return sessionKeyPrefix + hex.EncodeToString(digest[:])
}

func userSessionKey(uuid uuid.UUID) string {
	return fmt.Sprintf("user_session:%d", uuid)
}
