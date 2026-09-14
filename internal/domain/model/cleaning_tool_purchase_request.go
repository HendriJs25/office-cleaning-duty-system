package model

import "time"

type CleaningToolPurchaseRequest struct {
	ID                 int64
	OfficeID           int64
	CleaningToolID     *int64
	RequestedBy        int64
	ReviewedBy         *int64
	ItemName           string
	Quantity           int
	Reason             string
	Status             string
	EstimatedUnitPrice *int64
	ActualUnitPrice    *int64
	RejectionReason    *string
	SubmittedAt        *time.Time
	ApprovedAt         *time.Time
	RejectedAt         *time.Time
	PurchasedAt        *time.Time
	CompletedAt        *time.Time
	CreatedAt          time.Time
	UpdatedAt          *time.Time

	Office       Office        `gorm:"foreignKey:OfficeID;references:ID;"`
	CleaningTool *CleaningTool `gorm:"foreignKey:CleaningToolID;references:ID;"`
	Requester    User          `gorm:"foreignKey:RequestedBy;references:ID;constraint:OnDelete:RESTRICT;"`
	Reviewer     *User         `gorm:"foreignKey:ReviewedBy;references:ID;constraint:OnDelete:RESTRICT;"`
}

func (c *CleaningToolPurchaseRequest) TableName() string {
	return "cleaning_tool_purchase_requests"
}
