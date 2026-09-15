package repository

import (
	userrepository "cleaning/internal/repository/user"

	"gorm.io/gorm"
)

type Registry struct {
	User userrepository.Repository
}

func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{
		User: userrepository.NewRepository(db),
	}
}
