package cmd

import (
	"cleaning/internal/database"
	"cleaning/internal/database/seeders"
	"cleaning/internal/logger"
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

const seedTimeout = 30 * time.Second

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "run seeders",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSeeder()
	},
}

func runSeeder() error {
	postgresDB, err := database.NewPostgres(cfg.Database)
	if err != nil {
		return fmt.Errorf("connect to postgres failed: %w", err)
	}
	defer func() {
		if err := postgresDB.Close(); err != nil {
			logger.Log.WithError(err).Warn("failed to close postgres connection")
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), seedTimeout)
	defer cancel()

	logger.Log.Info("starting seeder")

	if err := seeders.Run(ctx, postgresDB.DB, cfg.Seed); err != nil {
		return err
	}

	logger.Log.Info("seeder finished")
	return nil
}
