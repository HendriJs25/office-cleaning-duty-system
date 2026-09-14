package model

import (
	"time"

	"gorm.io/gorm"
)

type CleaningTool struct {
	ID              int64
	OfficeID        int64
	Code            string
	Name            string
	Description     *string
	Quantity        int
	MinimumQuantity *int
	Unit            *string
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	DeletedAt       gorm.DeletedAt

	Office Office `gorm:"foreignKey:OfficeID;references:ID;constraint:OnDelete:RESTRICT;"`
}

func (c *CleaningTool) TableName() string {
	return "cleaning_tools"
}
