package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Employee struct {
	ID                  int64
	UUID                uuid.UUID
	OfficeID            int64
	FamilyName          string
	GivenName           string
	FamilyNameKana      *string
	GivenNameKana       *string
	IsActive            bool
	EmploymentStartDate *time.Time
	EmploymentEndDate   *time.Time
	CreatedAt           time.Time
	UpdatedAt           *time.Time
	DeletedAt           gorm.DeletedAt

	Office Office `gorm:"foreignKey:OfficeID;references:ID;constraint:OnDelete:RESTRICT;"`
}

func (e *Employee) TableName() string {
	return "employees"
}
