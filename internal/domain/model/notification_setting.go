package model

import "time"

type NotificationSetting struct {
	ID                      int64
	OfficeID                int64
	Provider                string
	EndpointURL             string
	Enabled                 bool
	NotifyScheduleConfirmed bool
	NotifyPurchaseSubmitted bool
	NotifyPurchaseApproved  bool
	NotifyPurchaseRejected  bool
	NotifyLowStock          bool
	CreatedAt               time.Time
	UpdatedAt               *time.Time

	Office Office `gorm:"foreignKey:OfficeID;references:ID"`
}

func (n *NotificationSetting) TableName() string {
	return "notification_settings"
}
