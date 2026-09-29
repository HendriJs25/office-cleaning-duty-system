package permissioncache

import (
	errConstant "cleaning/internal/constants/error"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	redislib "github.com/go-redis/redis/v8"
)

const (
	rolePermissionCacheKeyPrefix = "cache:role_permissions:"
	rolePermissionCacheTTL       = 30 * time.Minute
)

type repository struct {
	client redislib.Cmdable
}

type Repository interface {
	GetRolePermissions(context.Context, int64) ([]string, error)
	SetRolePermissions(context.Context, int64, []string) error
}

func NewRepository(client redislib.Cmdable) Repository {
	return &repository{
		client: client,
	}
}

func (r *repository) GetRolePermissions(ctx context.Context, roleID int64) ([]string, error) {
	payload, err := r.client.Get(ctx, rolePermissionCacheKey(roleID)).Bytes()
	if err != nil {
		if errors.Is(err, redislib.Nil) {
			return nil, errConstant.ErrNotFound
		}
		return nil, fmt.Errorf("get cache role permissions: %w", err)
	}

	var permissions []string
	if err = json.Unmarshal(payload, &permissions); err != nil {
		return nil, fmt.Errorf("unmarshal role permissions: %w", err)
	}
	return permissions, nil

}

func (r *repository) SetRolePermissions(ctx context.Context, roleID int64, permissions []string) error {
	payload, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("marshal cache role permissions: %w", err)
	}

	err = r.client.Set(ctx, rolePermissionCacheKey(roleID), payload, rolePermissionCacheTTL).Err()
	if err != nil {
		return fmt.Errorf("set cache role permissions: %w", err)
	}

	return nil
}

func rolePermissionCacheKey(roleID int64) string {
	return fmt.Sprintf("%s%d", rolePermissionCacheKeyPrefix, roleID)
}
