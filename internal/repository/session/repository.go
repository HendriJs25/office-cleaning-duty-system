package session

import (
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	Get(context.Context, string) (*model.Session, error)
	Delete(context.Context, string) error
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

func (r *repository) Get(ctx context.Context, accessToken string) (*model.Session, error) {
	payload, err := r.client.Get(ctx, sessionKey(accessToken)).Bytes()
	if err != nil {
		if errors.Is(err, redislib.Nil) {
			return nil, errConstant.ErrNotFound
		}
		return nil, fmt.Errorf("get session: %w", err)
	}

	var session model.Session
	if err := json.Unmarshal(payload, &session); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}
	return &session, nil
}

func (r *repository) Delete(ctx context.Context, accessToken string) error {
	deleted, err := r.client.Del(ctx, sessionKey(accessToken)).Result()
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	if deleted == 0 {
		return errConstant.ErrNotFound
	}
	return nil
}

func sessionKey(accessToken string) string {
	digest := sha256.Sum256([]byte(accessToken))

	return sessionKeyPrefix + hex.EncodeToString(digest[:])
}

func userSessionKey(uuid uuid.UUID) string {
	return fmt.Sprintf("user_session:%s", uuid.String())
}
