package db

import (
	"database/sql"
	"fmt"
	"log"
	"path"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/utils"
)

func NewMigrator(db *sql.DB, logger *zap.SugaredLogger) (*migrate.Migrate, error) {
	migrationsPath := path.Join(utils.GetProjectRoot(), "migrations")
	sourceUrl := fmt.Sprintf("file://%s", migrationsPath)
	if logger != nil {
		logger.Infof("Looking for migrations in: %s", sourceUrl)
	} else {
		fmt.Printf("Looking for migrations in: %s\n", sourceUrl)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create migrator driver: %s", err)
	}

	return migrate.NewWithDatabaseInstance(sourceUrl, "postgres", driver)
}
