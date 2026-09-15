package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	trimmedValue := strings.TrimSpace(value)

	if !exists || trimmedValue == "" {
		slog.Warn("use default value",
			"key", key,
			"fallback", fallback)
		return fallback
	}

	return trimmedValue
}

func getEnvRequired(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	trimmedValue := strings.TrimSpace(value)
	if !exists || trimmedValue == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return trimmedValue, nil
}

func getEnvInt(key string, fallback int) (int, error) {
	value, exists := os.LookupEnv(key)

	if !exists {
		return fallback, nil
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("failed to convert %s to int", value)
	}

	return valueInt, nil
}

func getEnvDuration(key string, fallback time.Duration) (time.Duration, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback, nil
	}

	result, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration:%w", key, err)
	}

	return result, nil
}

func validatePort(key, value string) error {
	port, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("%s must be a number: %w", key, err)
	}

	if port < 0 || port > 65535 {
		return fmt.Errorf("%s must be a number between 0 and 65535", key)
	}

	return nil
}
