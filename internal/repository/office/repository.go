package office

import (
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	"context"
	"errors"
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
	GetByID(context.Context, int64) (*model.Office, error)
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

func (r *repository) GetByID(ctx context.Context, id int64) (*model.Office, error) {
	var office model.Office
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&office).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrNotFound
		}
		return nil, fmt.Errorf("get office by id: %w", err)
	}
	return &office, nil
}
