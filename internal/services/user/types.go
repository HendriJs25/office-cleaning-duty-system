package user

import (
	"time"

	"github.com/google/uuid"
)

type GetUserResult struct {
	UserID      int64
	UUID        uuid.UUID
	UserName    string
	Email       string
	IsActive    bool
	LastLoginAt *time.Time
	RoleName    string
}
