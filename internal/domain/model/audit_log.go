package model

import (
	"time"

	"gorm.io/datatypes"
)

type AuditLog struct {
	ID         int64
	UserID     *int64
	Action     string
	EntityType string
	EntityID   *int64
	BeforeData *datatypes.JSON
	AfterData  *datatypes.JSON
	Metadata   *datatypes.JSON
	CreatedAt  time.Time

	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:RESTRICT;"`
}

func (a *AuditLog) TableName() string {
	return "audit_logs"
}
