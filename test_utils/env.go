package test_utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/michelm117/cycling-coach-lab/db"
)

type TestEnvironment struct {
	db *db.DatabaseEnv
}

func SetupEnvironment() TestEnvironment {
	// https://github.com/joho/godotenv/issues/43
	err := godotenv.Load(os.ExpandEnv("../.env"))
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	databaseEnv, err := db.GetDatabaseEnv()
	if err != nil {
		log.Fatalf("Error getting database environment: %s", err.Error())
	}
	return TestEnvironment{db: databaseEnv}
}
