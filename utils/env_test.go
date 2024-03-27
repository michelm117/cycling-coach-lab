package utils_test

import (
	"os"
	"testing"

	"github.com/michelm117/cycling-coach-lab/utils"
)

var ENVS = []string{"ENV", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "VERSION"}

func TestCheckForRequiredEnvVars(t *testing.T) {
	// Save the original environment variables
	originalEnv := make(map[string]string)
	for _, key := range ENVS {
		originalEnv[key] = os.Getenv(key)
	}

	// Clean up environment variables after the test
	defer func() {
		for key, value := range originalEnv {
			os.Setenv(key, value)
		}
	}()

	tests := []struct {
		desc    string
		envs    map[string]string
		wantErr bool
	}{
		{"Variable all set", map[string]string{
			"ENV":         "development",
			"DB_HOST":     "localhost",
			"DB_PORT":     "5432",
			"DB_USER":     "user",
			"DB_PASSWORD": "password",
			"DB_NAME":     "cycling_coach_lab",
			"VERSION":     "latest",
		}, false},

		{"Variables missing", map[string]string{
			"DB_HOST": "localhost",
			"DB_PORT": "5432",
			"DB_USER": "user",
			"DB_NAME": "cycling_coach_lab",
		}, true},
		// Add more test cases for other environment variables
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			// Clear the environment variables
			for _, key := range ENVS {
				os.Setenv(key, "")
			}
			// Set the environment variable for the test case
			for key, value := range tt.envs {
				os.Setenv(key, value)
			}
			// Run the function and check the error
			err := utils.CheckForRequiredEnvVars()
			if (err != nil) != tt.wantErr {
				t.Errorf("got error: %v, wantErr: %t for test %s", err, tt.wantErr, tt.desc)
			}
		})
	}
}
