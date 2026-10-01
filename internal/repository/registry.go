package repository

import (
	officerepository "cleaning/internal/repository/office"
	permissionrepository "cleaning/internal/repository/permission"
	rolerepository "cleaning/internal/repository/role"
	userrepository "cleaning/internal/repository/user"

	"gorm.io/gorm"
)

type Registry struct {
	User       userrepository.Repository
	Role       rolerepository.Repository
	Permission permissionrepository.Repository
	Office     officerepository.Repository
}

func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{
		User:       userrepository.NewRepository(db),
		Role:       rolerepository.NewRepository(db),
		Permission: permissionrepository.NewRepository(db),
		Office:     officerepository.NewRepository(db),
	}
}
