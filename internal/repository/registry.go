package repository

import (
	rolerepository "cleaning/internal/repository/role"
	userrepository "cleaning/internal/repository/user"

	"gorm.io/gorm"
)

type Registry struct {
	User userrepository.Repository
	Role rolerepository.Repository
}

func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{
		User: userrepository.NewRepository(db),
		Role: rolerepository.NewRepository(db),
	}
}
