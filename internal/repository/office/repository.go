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
	GetAllActive(context.Context) ([]model.Office, error)
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

func (r *repository) GetAllActive(ctx context.Context) ([]model.Office, error) {
	var offices []model.Office
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("id ASC").Find(&offices).Error
	if err != nil {
		return nil, fmt.Errorf("get all active offices: %w", err)
	}
	return offices, nil
}
