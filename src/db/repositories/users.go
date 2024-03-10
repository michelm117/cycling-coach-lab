package repositories

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/models"
)

type UsersRepository struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewUsersRepository(db *sql.DB, logger *zap.SugaredLogger) *UsersRepository {
	return &UsersRepository{
		db:     db,
		logger: logger,
	}
}

func (repo *UsersRepository) SearchForUser(keyword string) ([]*models.User, error) {
	rowName, err := repo.db.Query(
		"SELECT username, email FROM users WHERE users.username LIKE '%' || $1 || '%'  OR email LIKE '%' || $1 || '%'",
		keyword,
	)
	var users []*models.User
	for rowName.Next() {
		var user models.User
		err := rowName.Scan(&user.Name, &user.Email)
		if err != nil {
			return nil, fmt.Errorf("Error while trying to execute query: %s", err)
		}
		users = append(users, &user)
	}

	return users, err
}

func (repo *UsersRepository) GetById(id int) (*models.User, error) {
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

func (repo *UsersRepository) DeleteUser(user models.User) (*models.User, error) {
	_, err := repo.db.Exec(
		"DELETE FROM users WHERE users.username = $1 AND users.email = $2",
		user.Name,
		user.Email,
	)

	if err != nil {
		return nil, fmt.Errorf("Error while trying to execute query: %s", err)
	}

	return &user, nil
}

func (repo *UsersRepository) GetByName(name string) (*models.User, error) {
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

func (repo *UsersRepository) GetByEmail(email string) (*models.User, error) {
	row := repo.db.QueryRow("SELECT username, email FROM users WHERE users.email = $1", email)

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

func (repo *UsersRepository) AddUser(user models.User) (*models.User, error) {
	_, err := repo.db.Exec(
		"INSERT INTO users (username, email) VALUES ($1, $2)",
		user.Name,
		user.Email,
	)

	if err != nil {
		return nil, fmt.Errorf("User could no be added: %s", err)
	}

	return &user, nil
}

func (repo *UsersRepository) GetAllUsers() ([]*models.User, error) {
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

func (repo *UsersRepository) Count() (int, error) {
	row := repo.db.QueryRow("SELECT count(*) FROM users")
	var count int
	err := row.Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("Error while trying to execute query: %s", err)
	}
	return count, nil
}