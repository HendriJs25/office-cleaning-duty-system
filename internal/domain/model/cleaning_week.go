package model

import "time"

type CleaningWeek struct {
	ID            int64
	OfficeID      int64
	WeekStartDate time.Time
	WeekEndDate   time.Time
	Status        string
	CreatedBy     *int64
	ConfirmedBy   *int64
	ConfirmedAt   *time.Time
	CompletedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     *time.Time

	Office    Office `gorm:"foreignKey:OfficeID;references:ID;constraint:OnDelete:RESTRICT;"`
	Creator   *User  `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnDelete:RESTRICT;"`
	Confirmer *User  `gorm:"foreignKey:ConfirmedBy;references:ID;constraint:OnDelete:RESTRICT;"`
}

func (c *CleaningWeek) TableName() string {
	return "cleaning_weeks"
}
