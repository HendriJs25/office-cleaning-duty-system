package cmd

import (
	"cleaning/internal/database/migration"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "manage database migrations",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "run pending database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMigrateUp()
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "rollback last database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMigrateDown()
	},
}

var migrateForceCmd = &cobra.Command{
	Use:   "force [version]",
	Short: "force database migration version [-1 for no migrations]",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid migration version: %w", err)
		}
		return runMigrateForce(version)
	},
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the current database version",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMigrateVersion()
	},
}

func runMigrateUp() error {
	m, err := createMigrateInstance()
	if err != nil {
		return err
	}
	defer closeInstance(m)

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("nothing to migrate")
			return nil
		}
		return fmt.Errorf("migrate up: %w", err)
	}

	slog.Info("migration completed successfully")
	return nil
}

func runMigrateDown() error {
	m, err := createMigrateInstance()
	if err != nil {
		return err
	}
	defer closeInstance(m)

	if err = m.Steps(-1); err != nil {
		return fmt.Errorf("migrate down: %w", err)
	}

	slog.Info("migration rollback completed")
	return nil
}

func runMigrateForce(version int) error {
	m, err := createMigrateInstance()
	if err != nil {
		return err
	}
	defer closeInstance(m)

	if err = m.Force(version); err != nil {
		return fmt.Errorf("force migration version: %w", err)
	}

	slog.Info("migration version forced successfully",
		"version", version)
	return nil
}

func runMigrateVersion() error {
	m, err := createMigrateInstance()
	if err != nil {
		return err
	}
	defer closeInstance(m)

	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			slog.Info("no migrations has been applied")
			return nil
		}
		return fmt.Errorf("get migration version: %w", err)
	}

	slog.Info("current migration",
		"version", version,
		"dirty", dirty)

	return nil
}

func createMigrateInstance() (*migrate.Migrate, error) {
	m, err := migration.New(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("initialize migration: %w", err)
	}
	return m, nil
}

func closeInstance(m *migrate.Migrate) {
	sourceErr, databaseErr := m.Close()
	if sourceErr != nil {
		slog.Warn("failed to close migration source", "error", sourceErr)
	}

	if databaseErr != nil {
		slog.Warn("failed to close migration database", "error", databaseErr)
	}
}
