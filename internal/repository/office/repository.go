package office

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
	Create(context.Context, *model.Office) error
	ExistByCode(context.Context, string) (bool, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, office *model.Office) error {
	if err := r.db.WithContext(ctx).Create(office).Error; err != nil {
		return fmt.Errorf("create office: %w", err)
	}
	return nil
}

func (r *repository) ExistByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(model.Office{}).Where("code = ?", code).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("is office exist: %w", err)
	}
	return count > 0, nil
}
