package model

import "time"

type Role struct {
	ID          int64
	Code        string
	Name        string
	Description *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time

	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

func (r *Role) TableName() string {
	return "roles"
}
