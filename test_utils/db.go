package test_utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// https://medium.com/@dilshataliev/integration-tests-with-golang-test-containers-and-postgres-abb49e8096c5
type TestDatabase struct {
	Db        *sql.DB
	DbUrl     string
	container testcontainers.Container
}

func CreateTestContainer(
	ctx context.Context,
) (testcontainers.Container, *sql.DB, error) {
	globalEnv := SetupEnvironment()

	var env = map[string]string{
		"POSTGRES_USER":     globalEnv.databaseEnv.User,
		"POSTGRES_PASSWORD": globalEnv.databaseEnv.Password,
		"POSTGRES_DB":       globalEnv.databaseEnv.Name,
	}
	var port = "5432/tcp"
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16.2-alpine",
			ExposedPorts: []string{port},
			Cmd:          []string{"postgres", "-c", "fsync=off"},
			Env:          env,
			WaitingFor:   wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, fmt.Errorf("failed to start container: %s", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return container, nil, fmt.Errorf("failed to get container external port: %s", err)
	}

	log.Println("postgres container ready and running at port: ", mappedPort)

	url := fmt.Sprintf(
		"postgres://postgres:password@localhost:%s/%s?sslmode=disable",
		mappedPort.Port(),
		globalEnv.databaseEnv.Port,
	)
	db, err := sql.Open("postgres", url)
	if err != nil {
		return container, db, fmt.Errorf("failed to establish database connection: %s", err)
	}

	return container, db, nil
}

func NewPgMigrator(db *sql.DB) (*migrate.Migrate, error) {
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
