package model

import "time"

type VerificationToken struct {
	ID         int64
	EmployeeID int64
	HashToken  string
	TokenType  string
	ExpiresAt  time.Time
	CreatedAt  time.Time

	Employee Employee `gorm:"foreignKey:EmployeeID;references:ID;constraint:OnDelete:CASCADE"`
}
