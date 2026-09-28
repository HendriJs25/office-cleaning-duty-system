package permission

import (
	"cleaning/internal/domain/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindCodesByRoleID(context.Context, int64) ([]string, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindCodesByRoleID(ctx context.Context, roleID int64) ([]string, error) {
	var codes []string

	err := r.db.WithContext(ctx).Model(model.Permission{}).Joins(
		"JOIN role_permissions ON role_permissions.permission_id = permissions.id",
	).Where("role_permissions.role_id = ?", roleID).Pluck("permissions.code", &codes).Error

	if err != nil {
		return nil, fmt.Errorf("get permission codes: %w", err)
	}

	return codes, nil
}
