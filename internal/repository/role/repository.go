package role

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
	FindAll(context.Context) ([]model.Role, error)
	ExistByID(context.Context, int64) (bool, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role

	err := r.db.WithContext(ctx).Order("id ASC").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *repository) ExistByID(ctx context.Context, roleID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(model.Role{}).Where("id = ?", roleID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("exist role by id: %w", err)
	}
	return count > 0, nil
}
