package model

import "time"

type WeeklyAssignment struct {
	ID             int64
	CleaningWeekID int64
	EmployeeID     int64
	CleaningTaskID int64
	CreatedAt      time.Time
	UpdatedAt      *time.Time

	CleaningWeek CleaningWeek `gorm:"foreignKey:CleaningWeekID;references:ID;constraint:OnDelete:RESTRICT;"`
	Employee     Employee     `gorm:"foreignKey:EmployeeID;references:ID;constraint:OnDelete:RESTRICT;"`
	CleaningTask CleaningTask `gorm:"foreignKey:CleaningTaskID;references:ID;constraint:OnDelete:RESTRICT;"`
}

func (w *WeeklyAssignment) TableName() string {
	return "weekly_assignments"
}
