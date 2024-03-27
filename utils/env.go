package utils

import (
	"fmt"
	"os"
)

func CheckForRequiredEnvVars() error {
	requiredVars := []string{"ENV", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "VERSION"}

	for _, key := range requiredVars {
		value := os.Getenv(key)
		if value == "" {
			return fmt.Errorf("`%s` environment variable is required", key)
		}
	}
	return nil
}
