package repositories_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/michelm117/cycling-coach-lab/models"
	"github.com/michelm117/cycling-coach-lab/repositories"
	"github.com/michelm117/cycling-coach-lab/test_utils"
)

var DB *sql.DB

func TestMain(m *testing.M) {

	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16.2-alpine"),
		// postgres.WithInitScripts(filepath.Join("testdata", "init-user-db.sh")),
		// postgres.WithConfigFile(filepath.Join("testdata", "my-postgres.conf")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// Clean up the container
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// mappedPort, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	// if err != nil {
	// 	log.Fatalf("failed to get container external port: %s", err)
	// }

	dbURL, err := postgresContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to get container connection string: %s", err)
	}

	DB, err = sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("failed to open database: %s", err)
	}

	migrator, err := test_utils.NewPgMigrator(DB)
	if err != nil {
		log.Fatalf("failed to create migrator: %s", err)
	}
	if err := migrator.Up(); err != nil {
		log.Fatalf("failed to run migrations: %s", err)
	}
	// print user table content
	rows, err := DB.Query("SELECT * FROM users")
	if err != nil {
		log.Fatalf("failed to run query: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var email string
		err = rows.Scan(&id, &name, &email)
		if err != nil {
			log.Fatalf("failed to scan row: %s", err)
		}
		fmt.Printf("id: %d, name: %s, email: %s\n", id, name, email)

	}
	if err = rows.Err(); err != nil {
		log.Fatalf("failed to scan row: %s", err)
	}

	// Run the actual tests
	exitCode := m.Run()

	// Perform teardown tasks here

	// Exit with the exit code from the tests
	os.Exit(exitCode)

}

func TestCountUsers(t *testing.T) {
	repo := repositories.NewUserRepository(DB, nil)
	count, err := repo.Count()
	if err != nil {
		t.Errorf("Error while trying to count users: %s", err)
	}
	if count == 0 {
		t.Errorf("No users found")
	}
}

func TestAddUser(t *testing.T) {
	repo := repositories.NewUserRepository(DB, nil)
	beforeSize, err := repo.Count()
	if err != nil {
		t.Errorf("Error while trying to count users: %s", err)
	}
	u := models.User{
		Name:  "test",
		Email: "test@test.de",
	}
	user, err := repo.AddUser(u)
	if err != nil {
		t.Errorf("Error while trying to add a new user: %s", err)
	}

	if user == nil {
		t.Errorf("Newly added user was not returned: %s", u)
	}

	afterSize, err := repo.Count()
	if err != nil {
		t.Errorf("Error while trying to count users: %s", err)
	}
	if beforeSize+1 != afterSize {
		t.Errorf("Expected %d users, but got %d", beforeSize+1, afterSize)
	}
}

func TestGetByName(t *testing.T) {
	repo := repositories.NewUserRepository(DB, nil)
	user, err := repo.GetByName("user1")
	if err != nil {
		t.Errorf("Error while trying to get user by name: %s", err)
	}
	if user == nil {
		t.Errorf("User not found")
	}
}

func TestUserWithNameNotFound(t *testing.T) {
	repo := repositories.NewUserRepository(DB, nil)
	user, err := repo.GetByName("foo")
	if user != nil {
		t.Errorf("User should not be found")
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}
