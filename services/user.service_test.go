package services_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/assert"

	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/test_utils"
)

var DB *sql.DB

func getTestUser(userService services.UserServicer, t *testing.T) *model.User {
	// Search for user to get his id
	user, err := userService.GetByEmail("admin@example.com")
	if err != nil {
		t.Errorf("Error while trying to get user by name: %s", err)
	}
	return user
}

func TestMain(m *testing.M) {
	// Setup test environment
	ctx := context.Background()
	testDb := test_utils.CreateTestContainer(ctx)
	container := testDb.Container

	DB = testDb.DB

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
	userService := services.NewUserServicer(DB, nil)

	id := getTestUser(userService, t).ID
	user, err := userService.GetById(id)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestGetByEmail(t *testing.T) {
	repo := services.NewUserServicer(DB, nil)
	user, err := repo.GetByEmail("test@test.de")
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserWithEmailNotFound(t *testing.T) {
	repo := services.NewUserServicer(DB, nil)
	user, err := repo.GetByEmail("foo")
	assert.Nil(t, user)
	assert.NotNil(t, err)
}

func TestGetAllUsers(t *testing.T) {
	repo := services.NewUserServicer(DB, nil)
	expectedSize, err := repo.Count()
	assert.NoError(t, err)

	users, err := repo.GetAllUsers()
	assert.NoError(t, err)

	actualSize := len(users)
	assert.Equal(t, expectedSize, actualSize)
}

func TestAddUser(t *testing.T) {
	repo := services.NewUserServicer(DB, nil)
	beforeSize, err := repo.Count()
	assert.NoError(t, err)

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
	assert.NoError(t, err)
	assert.NotNil(t, user)

	afterSize, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, beforeSize+1, afterSize)
}

func TestDeleteUser(t *testing.T) {
	userService := services.NewUserServicer(DB, nil)
	expectedSize, err := userService.Count()
	expectedSize = expectedSize - 1
	assert.NoError(t, err)

	id := getTestUser(userService, t).ID
	err = userService.DeleteUser(id)
	assert.NoError(t, err)

	actualSize, err := userService.Count()
	assert.NoError(t, err)
	assert.Equal(t, expectedSize, actualSize)

	err = userService.DeleteUser(id)
	assert.NoError(t, err)
}

func TestCountUsers(t *testing.T) {
	repo := services.NewUserServicer(DB, nil)
	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Greater(t, count, 0)
}
