package repositories

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/models"
)

type UserRepository struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewUserRepository(db *sql.DB, logger *zap.SugaredLogger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (repo *UserRepository) GetById(id int) (*models.User, error) {
	row := repo.db.QueryRow("SELECT username, email FROM users WHERE users.username = $1", id)

	var user models.User
	err := row.Scan(&user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User with id '%d' not found", id)
		}
		return nil, fmt.Errorf("Error while trying to execute query: %s", err)
	}

	return &user, nil
}

func (repo *UserRepository) DeleteUser(user models.User) (*models.User, error) {
	println("AYPPP")
	println(user.Email)
	println(user.Name)
	_, err := repo.db.Exec("DELETE FROM users WHERE users.username = $1 AND users.email = $2", user.Name, user.Email)

	if err != nil {
		println("WE GET AN ERROR")
		println(err.Error())
		return nil, fmt.Errorf("Error while trying to execute query: %s", err)
	}

	return &user, nil
}

func (repo *UserRepository) GetByName(name string) (*models.User, error) {
	row := repo.db.QueryRow("SELECT username, email FROM users WHERE users.username = $1", name)

	var user models.User
	err := row.Scan(&user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User not found")
		}
		return nil, fmt.Errorf("Error while trying to execute query: %s", err)
	}
	return &user, nil
}

func (repo *UserRepository) AddUser(user models.User) (*models.User, error) {
	println("we are in add user" + user.Email)

	_, err := repo.db.Exec("INSERT INTO users (username, email) VALUES ($1, $2)", user.Name, user.Email)

	if err != nil {
		return nil, fmt.Errorf("User could no be added: %s", err)
	}

	return &user, nil
}

func (repo *UserRepository) GetAllUsers() ([]*models.User, error) {
	rows, err := repo.db.Query("SELECT username, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("Error while trying to execute query: %s", err)
	}

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Name, &user.Email)
		if err != nil {
			return nil, fmt.Errorf("Error while trying to execute query: %s", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error while trying to execute query: %s", err)
	}

	defer rows.Close()
	return users, nil
}
