package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           int64
	UUID         uuid.UUID
	RoleID       int64
	EmployeeID   *int64
	UserName     string
	Email        string
	PasswordHash string
	IsActive     bool
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    gorm.DeletedAt

	Role     Role      `gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:RESTRICT;"`
	Employee *Employee `gorm:"foreignKey:EmployeeID;references:ID;constraint:OnDelete:RESTRICT;"`
}

func (u *User) TableName() string {
	return "users"
}
