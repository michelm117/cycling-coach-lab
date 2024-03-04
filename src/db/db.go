package db

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/jackc/pgx/stdlib"
	"go.uber.org/zap"
)

func ConnectToDatabase(logger *zap.SugaredLogger) *sql.DB {
	logger.Info("Reading database environment variables")
	env, err := GetDatabaseEnv()
	if err != nil {
		log.Fatal("Error:", err)
	}

	logger.Info("Connecting to database")
	database, err := sql.Open("pgx", env.Address)
	if err != nil {
		log.Fatal("Error while connecting to db cause: " + err.Error())
	}

	if err := database.Ping(); err != nil {
		log.Fatalf("Error while pinging to db cause: %s", err.Error())
	}

	logger.Info("Creating migrator")
	migrator, err := NewMigrator(database, logger)
	if err != nil {
		log.Fatalf("failed to create migrator: %s", err)
	}

	logger.Info("Running migrations")
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrations: %s", err)
	}

	return database
}
