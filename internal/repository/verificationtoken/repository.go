package verificationtoken

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
	Create(context.Context, *model.VerificationToken) error
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

func (r *repository) DeleteByEmployeeIDAndTokenType(ctx context.Context, employeeID int64, tokenType string) error {
	err := r.db.WithContext(ctx).Where("employee_id = ? AND token_type = ?", employeeID, tokenType).Delete(&model.VerificationToken{}).Error
	if err != nil {
		return fmt.Errorf("delete token by employee_id and token_type: %w", err)
	}
	return nil
}
