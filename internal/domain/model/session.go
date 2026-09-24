package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UUID      uuid.UUID `json:"uuid"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	RoleCode  string    `json:"role_code"`
	CreatedAt time.Time `json:"created_at"`
}
