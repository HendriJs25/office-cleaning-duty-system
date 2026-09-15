package cmd

import (
	"cleaning/internal/database"
	"cleaning/internal/database/seeders"
	"context"
	"fmt"
	"log/slog"
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
			slog.Warn("failed to close postgres connection",
				"error", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), seedTimeout)
	defer cancel()

	if err := seeders.Run(ctx, postgresDB.DB, cfg.Seed); err != nil {
		return err
	}

	slog.Info("seed initialized successfully")
	return nil
}
