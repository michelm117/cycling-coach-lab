package db

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewMigrator(db *sql.DB) (*migrate.Migrate, error) {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("failed to get path")
	}

	projectRoot := filepath.Dir(filepath.Dir(currentFilePath))
	sourceUrl := fmt.Sprintf("file://%s/migrations", projectRoot)

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		log.Fatalf("failed to create migrator driver: %s", err)
	}

	return migrate.NewWithDatabaseInstance(sourceUrl, "postgres", driver)
}
