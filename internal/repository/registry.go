package repository

import (
	employeerepository "cleaning/internal/repository/employee"
	officerepository "cleaning/internal/repository/office"
	officearearepository "cleaning/internal/repository/officearea"
	permissionrepository "cleaning/internal/repository/permission"
	rolerepository "cleaning/internal/repository/role"
	userrepository "cleaning/internal/repository/user"
	verificationtokenrepository "cleaning/internal/repository/verificationtoken"

	"gorm.io/gorm"
)

type Registry struct {
	User              userrepository.Repository
	Role              rolerepository.Repository
	Permission        permissionrepository.Repository
	Office            officerepository.Repository
	Employee          employeerepository.Repository
	VerificationToken verificationtokenrepository.Repository
	OfficeArea        officearearepository.Repository
	Transaction       TransactionManager
}

func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{
		User:              userrepository.NewRepository(db),
		Role:              rolerepository.NewRepository(db),
		Permission:        permissionrepository.NewRepository(db),
		Office:            officerepository.NewRepository(db),
		Employee:          employeerepository.NewRepository(db),
		VerificationToken: verificationtokenrepository.NewRepository(db),
		OfficeArea:        officearearepository.NewRepository(db),
		Transaction:       NewTransactionManager(db),
	}
}
