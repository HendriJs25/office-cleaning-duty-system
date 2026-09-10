package config

import "fmt"

const (
	defaultDatabaseHost                   = "localhost"
	defaultDatabasePort                   = "5432"
	defaultDatabaseSSLMode                = "disable"
	defaultDatabaseMaxOpenConnections     = 10
	defaultDatabaseMaxLifetimeConnections = 10
	defaultDatabaseMaxIdleConnections     = 10
	defaultDatabaseMaxIdleTime            = 10
)

type Database struct {
	Host                   string
	Port                   string
	Name                   string
	Username               string
	Password               string
	SSLMode                string
	MaxOpenConnections     int
	MaxLifetimeConnections int
	MaxIdleConnections     int
	MaxIdleTime            int
}

func loadDatabase() (*Database, error) {
	name, err := getEnvRequired("DATABASE_NAME")
	if err != nil {
		return nil, err
	}

	username, err := getEnvRequired("DATABASE_USERNAME")
	if err != nil {
		return nil, err
	}

	password, err := getEnvRequired("DATABASE_PASSWORD")
	if err != nil {
		return nil, err
	}

	maxOpenConnections, err := getEnvInt("DATABASE_MAX_OPEN_CONNECTION", defaultDatabaseMaxOpenConnections)
	if err != nil {
		return nil, err
	}

	maxLifetimeConnections, err := getEnvInt("DATABASE_MAX_LIFETIME_CONNECTION", defaultDatabaseMaxLifetimeConnections)
	if err != nil {
		return nil, err
	}

	maxIdleConnections, err := getEnvInt("DATABASE_MAX_IDLE_CONNECTION", defaultDatabaseMaxIdleConnections)
	if err != nil {
		return nil, err
	}

	maxIdleTime, err := getEnvInt("DATABASE_MAX_IDLE_TIME", defaultDatabaseMaxIdleTime)
	if err != nil {
		return nil, err
	}

	return &Database{
		Host:                   getEnv("DATABASE_HOST", defaultDatabaseHost),
		Port:                   getEnv("DATABASE_PORT", defaultDatabasePort),
		Name:                   name,
		Username:               username,
		Password:               password,
		SSLMode:                getEnv("DATABASE_SSL_MODE", defaultDatabaseSSLMode),
		MaxOpenConnections:     maxOpenConnections,
		MaxLifetimeConnections: maxLifetimeConnections,
		MaxIdleConnections:     maxIdleConnections,
		MaxIdleTime:            maxIdleTime,
	}, nil
}

func (d *Database) Validate() error {
	if err := validatePort("DATABASE_PORT", d.Port); err != nil {
		return err
	}

	if d.MaxOpenConnections <= 0 {
		return fmt.Errorf("DATABASE_MAX_OPEN_CONNECTION must be greater than zero")
	}

	if d.MaxLifetimeConnections <= 0 {
		return fmt.Errorf("DATABASE_MAX_LIFETIME_CONNECTION must be greater than zero")
	}

	if d.MaxIdleConnections <= 0 {
		return fmt.Errorf("DATABASE_MAX_IDLE_CONNECTION must be greater than zero")
	}

	if d.MaxIdleTime <= 0 {
		return fmt.Errorf("DATABASE_MAX_IDLE_TIME must be greater than zero")
	}

	return nil
}
