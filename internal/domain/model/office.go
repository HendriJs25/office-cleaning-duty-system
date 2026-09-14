package model

import (
	"time"

	"gorm.io/gorm"
)

type Office struct {
	ID        int64
	Code      string
	Name      string
	Address   *string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt
}

func (o *Office) TableName() string {
	return "offices"
}
