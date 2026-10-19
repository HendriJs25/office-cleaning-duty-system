package role

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
	Create(context.Context, *model.Role) error
	AssignPermission(context.Context, int64, []model.Permission) error
	FindAll(context.Context) ([]model.Role, error)
	FindAllActive(context.Context) ([]model.Role, error)
	FindByID(context.Context, int64, bool) (*model.Role, error)
	FindByCode(context.Context, string) (*model.Role, error)
	ExistByCode(context.Context, string) (bool, error)
	UpdateStatus(context.Context, int64, bool) error
	HasAssignedUsers(context.Context, int64) (bool, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, role *model.Role) error {
	err := r.db.WithContext(ctx).Create(role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errConstant.ErrAlreadyExists
		}
		return fmt.Errorf("create role: %w", err)
	}
	return nil
}

func (r *repository) AssignPermission(ctx context.Context, roleID int64, permissions []model.Permission) error {
	role := &model.Role{
		ID: roleID,
	}
	err := r.db.WithContext(ctx).Model(role).Association("Permissions").Append(&permissions)
	if err != nil {
		return fmt.Errorf("assign permissions: %w", err)
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role

	err := r.db.WithContext(ctx).Order("id ASC").Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("find all roles: %w", err)
	}
	return roles, nil
}

func (r *repository) FindByID(ctx context.Context, id int64, withPermission bool) (*model.Role, error) {
	var role model.Role
	query := r.db.WithContext(ctx).Where("id = ?", id)

	if withPermission {
		query = query.Preload("Permissions")
	}

	err := query.Take(&role).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrRoleNotFound
		}
		return nil, fmt.Errorf("find role by id: %w", err)
	}
	return &role, nil
}

func (r *repository) FindAllActive(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role

	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("id ASC").Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("find all active roles: %w", err)
	}
	return roles, nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("code = ?", code).Take(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrNotFound
		}
		return nil, fmt.Errorf("find role by code: %w", err)
	}
	return &role, nil
}

func (r *repository) ExistByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Role{}).Where("code = ?", code).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("is exist by code: %w", err)
	}
	return count > 0, nil
}

func (r *repository) UpdateStatus(ctx context.Context, roleID int64, isActive bool) error {
	result := r.db.WithContext(ctx).Model(&model.Role{}).Where("id = ?", roleID).Update("is_active", isActive)
	if result.Error != nil {
		return fmt.Errorf("update role active status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errConstant.ErrRoleNotFound
	}

	return nil
}

func (r *repository) HasAssignedUsers(ctx context.Context, roleID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("role_id = ?", roleID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check assigned users: %w", err)
	}
	return count > 0, nil
}
