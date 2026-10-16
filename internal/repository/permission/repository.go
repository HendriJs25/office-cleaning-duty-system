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
	FindAll(context.Context) ([]model.Permission, error)
	FindCodesByRoleID(context.Context, int64) ([]string, error)
	FindByCodes(context.Context, []string) ([]model.Permission, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.WithContext(ctx).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
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

func (r *repository) FindByCodes(ctx context.Context, codes []string) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.WithContext(ctx).Where("code IN ?", codes).Find(&permissions).Error
	if err != nil {
		return nil, fmt.Errorf("get permissions by code: %w", err)
	}
	return permissions, nil
}
