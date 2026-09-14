package model

import "time"

type Permission struct {
	ID          int64
	Code        string
	Name        string
	Description *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func (p *Permission) TableName() string {
	return "permissions"
}
