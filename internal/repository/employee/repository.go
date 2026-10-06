package employee

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
	Create(context.Context, *model.Employee) error
	GetAllActive(context.Context) ([]model.Employee, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, employee *model.Employee) error {
	err := r.db.WithContext(ctx).Create(employee).Error
	if err != nil {
		return fmt.Errorf("create employee: %w", err)
	}
	return nil
}

func (r *repository) GetAllActive(ctx context.Context) ([]model.Employee, error) {
	var employees []model.Employee
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("id ASC").Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("get all active employees: %w", err)
	}
	return employees, nil
}
