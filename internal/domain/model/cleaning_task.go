package model

import (
	"time"

	"gorm.io/gorm"
)

type CleaningTask struct {
	ID           int64
	OfficeAreaID int64
	Code         string
	Name         string
	Description  *string
	SortOrder    int
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    gorm.DeletedAt

	OfficeArea OfficeArea `gorm:"foreignKey:OfficeAreaID;references:ID;constraint:OnDelete:RESTRICT;"`
}

func (c *CleaningTask) TableName() string {
	return "cleaning_tasks"
}
