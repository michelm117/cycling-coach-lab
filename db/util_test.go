package db_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/michelm117/cycling-coach-lab/db"
)

func prepareEnv(host, port, user, password, dbname string) {
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port)
	os.Setenv("DB_USER", user)
	os.Setenv("DB_PASSWORD", password)
	os.Setenv("DB_NAME", dbname)
}

func TestGetdatabaseEnv(t *testing.T) {
	tests := []struct {
		name          string
		missingEnvVar string
	}{
		{"AllEnvSet", ""},
		{"MissingHost", "DB_HOST"},
		{"MissingPort", "DB_PORT"},
		{"MissingUser", "DB_USER"},
		{"MissingPassword", "DB_PASSWORD"},
		{"MissingDBName", "DB_NAME"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prepareEnv("localhost", "5432", "postgres", "postgres", "postgres")
			// Unset the environment variable if needed
			if tt.missingEnvVar != "" {
				os.Unsetenv(tt.missingEnvVar)
			}

			got, err := db.GetDatabaseEnv()

			// Check if an error is expected
			if tt.missingEnvVar != "" {
				if err == nil {
					t.Errorf(
						"Expected an error for missing environment variable %s, but got nil",
						tt.missingEnvVar,
					)
				}
				return
			}

			// No error expected, compare the result
			expected := fmt.Sprintf(
				"postgres://%s:%s@localhost:5432/postgres",
				"postgres",
				"postgres",
			)
			if err != nil {
				t.Errorf("Unexpected error while building psql url: %s", err)
			}
			if got.Address != expected {
				t.Errorf("Expected %s but got %s", expected, got)
			}
		})
	}
}
