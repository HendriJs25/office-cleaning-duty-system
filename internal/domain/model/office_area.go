package model

import (
	"time"

	"gorm.io/gorm"
)

type OfficeArea struct {
	ID          int64
	OfficeID    int64
	Code        string
	Name        string
	Description *string
	SortOrder   int
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   gorm.DeletedAt

	Office Office `gorm:"foreignKey:OfficeID;references:ID;constraint:OnDelete:RESTRICT;"`
}

func (o *OfficeArea) TableName() string {
	return "office_areas"
}
