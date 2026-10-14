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
