package cmd

import (
	"cleaning/internal/common/validator"
	"cleaning/internal/database"
	"cleaning/internal/handler"
	"cleaning/internal/repository"
	"cleaning/internal/repository/session"
	"cleaning/internal/routes"
	"cleaning/internal/services"
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
	postgresDB, err := database.NewPostgres(cfg.Database)
	if err != nil {
		return fmt.Errorf("connect to postgres failed: %w", err)
	}

	defer func() {
		if err := postgresDB.Close(); err != nil {
			slog.Warn("failed to close postgres connection", "error", err)
		}
	}()

	slog.Info("postgres connection established",
		"host", cfg.Database.Host,
		"port", cfg.Database.Port,
		"name", cfg.Database.Name)

	redisDB, err := database.NewRedis(cfg.Redis)
	if err != nil {
		return fmt.Errorf("connect to redis failed: %w", err)
	}

	defer func() {
		if err := redisDB.Close(); err != nil {
			slog.Warn("failed to close redis connection", "error", err)
		}
	}()

	slog.Info("redis connection established",
		"host", cfg.Redis.Host,
		"port", cfg.Redis.Port,
		"database", cfg.Redis.DB)

	sessionRepository := session.NewRepository(redisDB.Client)
	repositoryRegistry := repository.NewRegistry(postgresDB.DB)
	serviceRegistry := services.NewRegistry(repositoryRegistry, sessionRepository, cfg.JWT)
	v := validator.New()
	handlerRegistry := handler.NewRegistry(serviceRegistry, v)

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
