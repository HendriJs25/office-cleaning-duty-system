package cmd

import (
	"cleaning/internal/handler"
	"cleaning/internal/routes"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServer()
	},
}

func runServer() error {
	handlerRegistry := handler.NewRegistry()

	router := gin.Default()
	group := router.Group("api/v1")
	routerRegistry := routes.NewRegistry(group, handlerRegistry)
	routerRegistry.Register()

	slog.Info("starting cleaning app",
		"environment", cfg.App.Env,
		"address", cfg.ServerAddress())

	if err := router.Run(cfg.ServerAddress()); err != nil {
		return fmt.Errorf("starting server failed: %w", err)
	}
	return nil
}
