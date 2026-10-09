package employee

import (
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	Create(context.Context, *model.Employee) error
	FindAllActive(context.Context) ([]model.Employee, error)
	FindByID(context.Context, *int64) (*model.Employee, error)
	FindAssignableEmployees(context.Context, *uuid.UUID) ([]model.Employee, error)
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

func (r *repository) FindAllActive(ctx context.Context) ([]model.Employee, error) {
	var employees []model.Employee
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("id ASC").Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("find all active employees: %w", err)
	}
	return employees, nil
}

func (r *repository) FindByID(ctx context.Context, id *int64) (*model.Employee, error) {
	var employee model.Employee
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&employee).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("find employee by id: %w", err)
	}
	return &employee, nil
}

func (r *repository) FindAssignableEmployees(ctx context.Context, userUUID *uuid.UUID) ([]model.Employee, error) {
	var employees []model.Employee

	subquery := r.db.WithContext(ctx).Table("users").Select("1").Where("users.employee_id = employees.id")

	if userUUID != nil {
		subquery = subquery.Where("users.uuid <> ?", *userUUID)
	}

	err := r.db.WithContext(ctx).Model(&model.Employee{}).Where("employees.is_active = ?", true).Where("NOT EXISTS (?)", subquery).Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("find all assignable employees: %w", err)
	}
	return employees, nil
}
