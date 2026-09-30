package user

import (
	errConstant "cleaning/internal/constants/error"
	"cleaning/internal/domain/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindAll(context.Context) ([]model.User, error)
	FindByEmail(context.Context, string) (*model.User, error)
	FindByUUID(context.Context, uuid.UUID) (*model.User, error)
	ExistByEmail(context.Context, string) (bool, error)
	Create(context.Context, *model.User) error
	UpdateLastLoginAt(context.Context, uuid.UUID) error
	UpdateStatus(context.Context, uuid.UUID, bool) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := r.db.WithContext(ctx).Joins("Role").Order("id ASC").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	return users, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &user, nil
}

func (r *repository) FindByUUID(ctx context.Context, uuid uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Joins("Role").Joins("Employee").Joins("Employee.Office").Where("users.uuid = ?", uuid).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrNotFound
		}
		return nil, fmt.Errorf("find user by uuid: %w", err)
	}
	return &user, nil
}

func (r *repository) ExistByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("is email exist: %w", err)
	}
	return count > 0, nil
}

func (r *repository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *repository) UpdateLastLoginAt(ctx context.Context, userUUID uuid.UUID) error {
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("uuid = ?", userUUID).Update("last_login_at", time.Now().UTC()).Error
	if err != nil {
		return fmt.Errorf("update last login at: %w", err)
	}
	return nil
}

func (r *repository) UpdateStatus(ctx context.Context, userUUID uuid.UUID, status bool) error {
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("uuid = ?", userUUID).Update("is_active", status).Error
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	return nil
}
