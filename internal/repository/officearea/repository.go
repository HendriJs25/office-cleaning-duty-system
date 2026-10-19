package officearea

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
	DeactivateByOfficeID(context.Context, int64) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) DeactivateByOfficeID(ctx context.Context, officeID int64) error {
	result := r.db.WithContext(ctx).Model(&model.OfficeArea{}).Where("office_id = ? AND is_active = ?", officeID, true).Update("is_active", false)
	if result.Error != nil {
		return fmt.Errorf("deactivate office area by office id failed: %w", result.Error)
	}
	return nil
}
