package repositories_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/michelm117/cycling-coach-lab/models"
	"github.com/michelm117/cycling-coach-lab/repositories"
	"github.com/michelm117/cycling-coach-lab/test_utils"
)

var DB *sql.DB

func TestMain(t *testing.M) {
	testDb := test_utils.SetupTestDatabase()
	DB = testDb.Db
	defer testDb.TearDown()
	os.Exit(t.Run())
}

func TestAddUser(t *testing.T) {
	repo := repositories.NewUserRepository(DB, nil)
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
