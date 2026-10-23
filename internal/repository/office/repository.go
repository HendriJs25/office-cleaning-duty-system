package office

import (
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	Create(context.Context, *model.Office) error
	ExistByCode(context.Context, string) (bool, error)
	GetAll(context.Context) ([]model.Office, error)
	GetAllActive(context.Context) ([]model.Office, error)
	GetByID(context.Context, int64) (*model.Office, error)
	GetByIDForUpdate(context.Context, int64) (*model.Office, error)
	UpdateByID(context.Context, int64, map[string]any) error
	UpdateStatus(context.Context, int64, bool) error
	HasAssignedEmployee(context.Context, int64) (bool, error)
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

func (r *repository) GetAll(ctx context.Context) ([]model.Office, error) {
	var offices []model.Office
	err := r.db.WithContext(ctx).Order("id ASC").Find(&offices).Error
	if err != nil {
		return nil, fmt.Errorf("get offices: %w", err)
	}
	return offices, nil
}

func (r *repository) GetAllActive(ctx context.Context) ([]model.Office, error) {
	var offices []model.Office
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("id ASC").Find(&offices).Error
	if err != nil {
		return nil, fmt.Errorf("get all active offices: %w", err)
	}
	return offices, nil
}

func (r *repository) GetByID(ctx context.Context, officeID int64) (*model.Office, error) {
	var office model.Office
	err := r.db.WithContext(ctx).Where("id = ?", officeID).Take(&office).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrNotFound
		}
		return nil, fmt.Errorf("get office by id: %w", err)
	}
	return &office, nil
}

func (r *repository) GetByIDForUpdate(ctx context.Context, officeID int64) (*model.Office, error) {
	var office model.Office
	err := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", officeID).Take(&office).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrNotFound
		}
		return nil, fmt.Errorf("get office by id: %w", err)
	}
	return &office, nil
}

func (r *repository) UpdateByID(ctx context.Context, officeID int64, fields map[string]any) error {
	result := r.db.WithContext(ctx).Model(model.Office{}).Where("id = ?", officeID).Updates(fields)
	if result.Error != nil {
		return fmt.Errorf("update office by id: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errConstant.ErrNotFound
	}
	return nil
}

func (r *repository) UpdateStatus(ctx context.Context, officeID int64, isActive bool) error {
	result := r.db.WithContext(ctx).Model(model.Office{}).Where("id = ?", officeID).Update("is_active", isActive)
	if result.Error != nil {
		return fmt.Errorf("update office active status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errConstant.ErrNotFound
	}
	return nil
}

func (r *repository) HasAssignedEmployee(ctx context.Context, officeID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(model.Employee{}).Where("office_id = ?", officeID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check assigned employee: %w", err)
	}
	return count > 0, nil
}
