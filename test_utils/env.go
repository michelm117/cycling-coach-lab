package test_utils

import (
	"log"
	"os"

	"github.com/michelm117/cycling-coach-lab/db"
)

type TestEnvironment struct {
	db *db.DatabaseEnv
}

func SetupEnvironment() TestEnvironment {
	// https://github.com/joho/godotenv/issues/43
	os.Setenv("PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "postgres")

	databaseEnv, err := db.GetDatabaseEnv()
	if err != nil {
		log.Fatalf("Error getting database environment: %s", err.Error())
	}
	return TestEnvironment{db: databaseEnv}
}
