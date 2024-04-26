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

type Migrator interface {
	Up() error
	Down() error
	Force(version int) error
	Reset() error
}

type Migrate struct {
	db       *sql.DB
	migrator *migrate.Migrate
	logger   *zap.SugaredLogger
}

func NewMigrator(
	db *sql.DB,
	migrationFolder string,
	logger *zap.SugaredLogger,
) Migrator {
	migrationsPath := path.Join(utils.GetProjectRoot(), migrationFolder)
	sourceUrl := fmt.Sprintf("file://%s", migrationsPath)
	logger.Infof("Looking for migrations in: %s", sourceUrl)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Panicf("failed to create postgres driver: %s", err)
	}
	migrator, err := migrate.NewWithDatabaseInstance(sourceUrl, "postgres", driver)
	if err != nil {
		log.Panicf("failed to create migrator: %s", err)
	}

	return &Migrate{
		db:       db,
		migrator: migrator,
		logger:   logger,
	}
}

func (m *Migrate) Up() error {
	m.logger.Info("Running migrations")
	if err := m.migrator.Up(); err != nil && err != migrate.ErrNoChange {
		m.logger.Errorf("Failed to run migrations: %w", err)
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func (m *Migrate) Down() error {
	m.logger.Info("Resetting migrations")
	if err := m.migrator.Down(); err != nil && err != migrate.ErrNoChange {
		m.logger.Errorf("Failed to reset migrations: %w", err)
		return fmt.Errorf("failed to reset migrations: %w", err)
	}
	return nil
}

func (m *Migrate) Force(version int) error {
	m.logger.Infof("Forcing migration to version: %d", version)
	if err := m.migrator.Force(version); err != nil {
		m.logger.Errorf("Failed to force migration: %w", err)
		return fmt.Errorf("failed to force migration: %w", err)
	}
	return nil
}

func (m *Migrate) Reset() error {
	if err := m.Down(); err != nil {
		m.logger.Error(err)
		return utils.Danger(err.Error())
	}
	if err := m.Up(); err != nil {
		m.logger.Error(err)
		return utils.Danger(err.Error())
	}
	return nil
}
