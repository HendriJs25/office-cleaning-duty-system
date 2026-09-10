package database

import (
	"cleaning/internal/config"
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const databasePingTimeout = 5 * time.Second

type Postgres struct {
	DB    *gorm.DB
	sqlDB *sql.DB
}

func NewPostgres(cfg *config.Database) (*Postgres, error) {
	gormDB, err := gorm.Open(postgres.Open(BuildPostgresURL(cfg)), &gorm.Config{
		Logger:               logger.Default.LogMode(logger.Info),
		DisableAutomaticPing: true,
		TranslateError:       true,
	})
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql database handle: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetimeConnections) * time.Second)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), databasePingTimeout)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Postgres{
		DB:    gormDB,
		sqlDB: sqlDB,
	}, nil
}

func (p *Postgres) Close() error {
	return p.sqlDB.Close()
}

func BuildPostgresURL(cfg *config.Database) string {
	connectionURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		Path:   "/" + cfg.Name,
	}

	query := connectionURL.Query()
	query.Set("sslmode", cfg.SSLMode)
	connectionURL.RawQuery = query.Encode()
	return connectionURL.String()
}
