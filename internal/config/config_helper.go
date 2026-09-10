package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
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
