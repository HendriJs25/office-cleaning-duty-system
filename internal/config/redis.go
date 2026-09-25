package config

import "fmt"

const (
	defaultRedisHost = "localhost"
	defaultRedisPort = "6379"
	defaultRedisDB   = 0
)

type Redis struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func loadRedis() (*Redis, error) {
	redisDB, err := getEnvInt("REDIS_DB", defaultRedisDB)
	if err != nil {
		return nil, err
	}

	return &Redis{
		Host:     getEnv("REDIS_HOST", defaultRedisHost),
		Port:     getEnv("REDIS_PORT", defaultRedisPort),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       redisDB,
	}, nil
}

func (r *Redis) Validate() error {
	if err := validatePort("REDIS_PORT", r.Port); err != nil {
		return err
	}

	if r.DB < 0 {
		return fmt.Errorf("REDIS_DB must not be negative")
	}
	return nil
}
