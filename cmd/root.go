package cmd

import (
	"cleaning/internal/config"
	"fmt"

	"github.com/spf13/cobra"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "cleaning-app",
	Short: "cleaning application",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("load config failed: %w", err)
		}
		return nil
	},
}

func Execute() error {
	registerCommands()

	return rootCmd.Execute()
}

func registerCommands() {
	rootCmd.AddCommand(
		serveCmd,
		migrateCmd,
		seedCmd)

	migrateCmd.AddCommand(
		migrateUpCmd,
		migrateDownCmd,
		migrateForceCmd,
		migrateVersionCmd)
}
