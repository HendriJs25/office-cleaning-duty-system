package verificationtoken

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
	Create(context.Context, *model.VerificationToken) error
	FindByEmployeeIDAndTokenType(context.Context, int64, string) (*model.VerificationToken, error)
	FindByTokenHashAndTokenType(context.Context, string, string) (*model.VerificationToken, error)
	DeleteByEmployeeIDAndTokenType(context.Context, int64, string) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, token *model.VerificationToken) error {
	err := r.db.WithContext(ctx).Create(token).Error
	if err != nil {
		return fmt.Errorf("create token: %w", err)
	}
	return nil
}

func (r *repository) FindByEmployeeIDAndTokenType(ctx context.Context, employeeID int64, tokenType string) (*model.VerificationToken, error) {
	var token model.VerificationToken
	err := r.db.WithContext(ctx).Where("employee_id = ? AND token_type = ?", employeeID, tokenType).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrNotFound
		}
		return nil, fmt.Errorf("find token by employee id: %w", err)
	}
	return &token, nil
}

func (r *repository) FindByTokenHashAndTokenType(ctx context.Context, tokenHash string, tokenType string) (*model.VerificationToken, error) {
	var token model.VerificationToken
	err := r.db.WithContext(ctx).Where("hash_token = ? AND token_type = ?", tokenHash, tokenType).Take(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrNotFound
		}
		return nil, fmt.Errorf("find token by token hash: %w", err)
	}
	return &token, nil
}

func (r *repository) DeleteByEmployeeIDAndTokenType(ctx context.Context, employeeID int64, tokenType string) error {
	err := r.db.WithContext(ctx).Where("employee_id = ? AND token_type = ?", employeeID, tokenType).Delete(&model.VerificationToken{}).Error
	if err != nil {
		return fmt.Errorf("delete token by employee_id and token_type: %w", err)
	}
	return nil
}
