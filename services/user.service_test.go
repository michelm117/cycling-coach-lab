package services_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/test_utils"
)

var DB *sql.DB

func TestMain(m *testing.M) {
	// Setup test environment
	ctx := context.Background()
	testDb := test_utils.CreateTestContainer(ctx)
	container := testDb.Container

	DB = testDb.Db

	// Run the actual tests
	exitCode := m.Run()

	// Perform tear down
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// Exit with the exit code from the tests
	os.Exit(exitCode)

}

func TestGetById(t *testing.T) {
	repo := services.NewUserService(DB, nil)
	// Search for user to get his id
	user, err := repo.GetByEmail("admin@example.com")
	if err != nil {
		t.Errorf("Error while trying to get user by name: %s", err)
	}

	id := user.ID
	user, err = repo.GetById(id)
	if err != nil {
		t.Errorf("Error while trying to get user by ID: %s", err)
	}

	if user == nil {
		t.Errorf("User not found")
	}
}

func TestGetByEmail(t *testing.T) {
	repo := services.NewUserService(DB, nil)
	user, err := repo.GetByEmail("test@test.de")
	if err != nil {
		t.Errorf("Error while trying to get user by email: %s", err)
	}
	if user == nil {
		t.Errorf("User not found")
	}
}

func TestUserWithEmailNotFound(t *testing.T) {
	repo := services.NewUserService(DB, nil)
	user, err := repo.GetByEmail("foo")
	if user != nil {
		t.Errorf("User should not be found")
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestGetAllUsers(t *testing.T) {
	repo := services.NewUserService(DB, nil)
	expectedSize, err := repo.Count()
	if err != nil {
		t.Errorf("Error while trying to count users: %s", err)
	}

	users, err := repo.GetAllUsers()
	if err != nil {
		t.Errorf("Error while trying to get all users: %s", err)

	}

	actualSize := len(users)
	if actualSize != expectedSize {
		t.Errorf("actual size %v is not expectedSize %v", actualSize, expectedSize)
	}
}

func TestAddUser(t *testing.T) {
	repo := services.NewUserService(DB, nil)
	beforeSize, err := repo.Count()
	if err != nil {
		t.Errorf("Error while trying to count users: %s", err)
	}
	u := model.User{
		Firstname:    "first",
		Lastname:     "last",
		Email:        "foo@bar.com",
		DateOfBirth:  time.Now(),
		Role:         "admin",
		Status:       "active",
		PasswordHash: "hash",
	}
	user, err := repo.AddUser(u)
	if err != nil {
		t.Errorf("Error while trying to add a new user: %s", err)
	}

	if user == nil {
		t.Errorf("Newly added user was not returned: %v", u)
	}

	afterSize, err := repo.Count()
	if err != nil {
		t.Errorf("Error while trying to count users: %s", err)
	}
	if beforeSize+1 != afterSize {
		t.Errorf("Expected %d users, but got %d", beforeSize+1, afterSize)
	}
}

func TestDeleteUser(t *testing.T) {
	repo := services.NewUserService(DB, nil)
	expectedSize, err := repo.Count()
	expectedSize--
	if err != nil {
		t.Errorf("Error while trying to count users: %s", err)
	}

	err = repo.DeleteUser("jan@ullrich.de")
	if err != nil {
		t.Errorf("Error while trying to delete a users: %s", err)
	}

	actualSize, err := repo.Count()
	if err != nil {
		t.Errorf("Error while trying to count users after deleting one: %s", err)
	}

	if actualSize != expectedSize {
		t.Errorf("actual size %v is not expectedSize %v", actualSize, expectedSize)
	}

	err = repo.DeleteUser("jan@ullrich.de")
	if err != nil {
		t.Errorf("Deleting an user that does not exists should not throw any errors: %s", err)
	}
}

func TestCountUsers(t *testing.T) {
	repo := services.NewUserService(DB, nil)
	count, err := repo.Count()
	if err != nil {
		t.Errorf("Error while trying to count users: %s", err)
	}
	if count == 0 {
		t.Errorf("No users found")
	}
}
