package test_utils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/testcontainers/testcontainers-go"
	testcontainerPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/michelm117/cycling-coach-lab/db"
)

// https://medium.com/@dilshataliev/integration-tests-with-golang-test-containers-and-postgres-abb49e8096c5
type TestDatabase struct {
	Db        *sql.DB
	DbUrl     string
	container testcontainers.Container
}

func SetupTestDatabase() *TestDatabase {
	env := SetupEnvironment().databaseEnv
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)

	// setup db container
	container, dbInstance, dbAddr, err := createContainer(ctx, env)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}

	// migrate db schema
	err = migrateDb(dbAddr)
	if err != nil {
		log.Fatal("failed to perform db migration", err)
	}
	cancel()

	return &TestDatabase{
		container: container,
		Db:        dbInstance,
		DbUrl:     dbAddr,
	}
}

func (tdb *TestDatabase) TearDown() {
	tdb.Db.Close()
	// remove test container
	_ = tdb.container.Terminate(context.Background())
}

func createContainer(
	ctx context.Context,
	env *db.DatabaseEnv,
) (testcontainers.Container, *sql.DB, string, error) {

	postgresContainer, err := testcontainerPostgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16.2-alpine"),
		testcontainerPostgres.WithDatabase(env.Name),
		testcontainerPostgres.WithUsername(env.User),
		testcontainerPostgres.WithPassword(env.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		err = fmt.Errorf("failed to start container: %s", err)
	}

	if err != nil {
		err = fmt.Errorf("Error: %s", err)
	}
	db, err := sql.Open("pgx", env.Address)
	if err != nil {
		err = fmt.Errorf("Error while connecting to db cause: " + err.Error())
	}

	if err := db.Ping(); err != nil {
		err = fmt.Errorf("Error while pinging to db cause: " + err.Error())
	}

	return postgresContainer, db, env.Address, err

}

func migrateDb(dbAddr string) error {
	m, err := migrate.New(
		fmt.Sprintf("file:///home/michelm/Projects/cycling-coach-lab/migrations"),
		fmt.Sprintf("%s?sslmode=disable", dbAddr),
	)
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("migration done")

	return nil
}
