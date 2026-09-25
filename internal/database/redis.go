package database

import (
	"cleaning/internal/config"
	"context"
	"fmt"
	"net"
	"time"

	redislib "github.com/go-redis/redis/v8"
)

const redisPingTimeout = time.Second * 5

type Redis struct {
	Client *redislib.Client
}

func NewRedis(cfg *config.Redis) (*Redis, error) {
	client := redislib.NewClient(&redislib.Options{
		Addr:     net.JoinHostPort(cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), redisPingTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &Redis{
		Client: client,
	}, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}
