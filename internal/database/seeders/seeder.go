package seeders

import (
	"cleaning/internal/config"
	"context"
	"fmt"

	"gorm.io/gorm"
)

func Run(ctx context.Context, db *gorm.DB, cfg *config.Seed) error {
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		permissions, err := seedPermissions(tx)
		if err != nil {
			return fmt.Errorf("seed permissions: %w", err)
		}

		roles, err := seedRoles(tx)
		if err != nil {
			return fmt.Errorf("seed roles: %w", err)
		}

		if err := seedRolePermission(tx, roles, permissions); err != nil {
			return fmt.Errorf("seed role permissions: %w", err)
		}

		if err := seedAdmin(tx, roles["admin"], cfg); err != nil {
			return fmt.Errorf("seed admin: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
