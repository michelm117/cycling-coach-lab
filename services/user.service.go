package services

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/model"
)

type UserService struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewUserService(db *sql.DB, logger *zap.SugaredLogger) *UserService {
	return &UserService{
		db:     db,
		logger: logger,
	}
}

func (repo *UserService) GetById(id int) (*model.User, error) {
	row := repo.db.QueryRow("SELECT * FROM users WHERE users.id = $1", id)

	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Firstname,
		&user.Lastname,
		&user.DateOfBirth,
		&user.PasswordHash,
		&user.Status,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id '%d' not found", id)
		}
		return nil, fmt.Errorf("error while trying to execute query: %s", err)
	}

	return &user, nil
}

func (repo *UserService) GetByEmail(email string) (*model.User, error) {
<<<<<<< HEAD
	row := repo.db.QueryRow("SELECT * FROM users WHERE users.email = $1", email)

	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Firstname,
		&user.Lastname,
		&user.DateOfBirth,
		&user.PasswordHash,
		&user.Status,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
=======
	row := repo.db.QueryRow("SELECT username, email, password, admin FROM users WHERE users.email = $1", email)

	var user model.User
	err := row.Scan(&user.Name, &user.Email, &user.Password, &user.Admin)
>>>>>>> 569d93b (auth)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User with email '%s' not found", email)
		}
		return nil, fmt.Errorf("error while trying to execute query: %s", err)
	}
	return &user, nil
}

<<<<<<< HEAD
=======
func (repo *UserService) AddUser(user model.User) (*model.User, error) {
	_, err := repo.db.Exec(
		"INSERT INTO users (username, email, password, admin) VALUES ($1, $2, $3, $4)",
		user.Name,
		user.Email,
		user.Password,
		user.Admin,
	)

	if err != nil {
		return nil, fmt.Errorf("user could no be added: %s", err)
	}

	return &user, nil
}

>>>>>>> 569d93b (auth)
func (repo *UserService) GetAllUsers() ([]*model.User, error) {
	rows, err := repo.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("error while trying to execute query: %s", err)
	}

	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Firstname,
			&user.Lastname,
			&user.DateOfBirth,
			&user.PasswordHash,
			&user.Status,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error while trying to execute query: %s", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while trying to execute query: %s", err)
	}

	defer rows.Close()
	return users, nil
}

func (repo *UserService) AddUser(user model.User) (*model.User, error) {
	_, err := repo.db.Exec(
		"INSERT INTO users (email, firstname, lastname, date_of_birth, password_hash, status, role) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.Email,
		user.Firstname,
		user.Lastname,
		user.DateOfBirth,
		user.PasswordHash,
		user.Status,
		user.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("user could not be added: %s", err)
	}

	return &user, nil
}

func (repo *UserService) DeleteUser(id int) error {
	_, err := repo.db.Exec(
		"DELETE FROM users WHERE users.id = $1",
		id,
	)

	if err != nil {
		return fmt.Errorf("error while trying to execute query: %s", err)
	}
	return nil
}

func (repo *UserService) Count() (int, error) {
	row := repo.db.QueryRow("SELECT count(*) FROM users")
	var count int
	err := row.Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("error while trying to execute query: %s", err)
	}
	return count, nil
}

func (repo *UserService) AddSessionId(email string, sessionId string) error {
	fmt.Println("error weee")
	_, err := repo.db.Exec(
		"UPDATE users SET session_id = $1 WHERE email = $2",
		sessionId,
		email,
	)
	fmt.Println("AddSessionId")
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserService) VeryifyUserSessionId(sessionId string) (bool, error) {
	row := repo.db.QueryRow("SELECT username, password, email FROM users WHERE users.session_id = $1", sessionId)
	var user model.User
	err := row.Scan(&user.Name, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("user not found")
		}
		return false, fmt.Errorf("error while trying to execute query: %s", err)
	}
	return true, nil
}

func (repo *UserService) GetUserBySessionId(sessionId string) (*model.User, error) {
	row := repo.db.QueryRow("SELECT username, password, email FROM users WHERE users.session_id = $1", sessionId)
	var user model.User
	err := row.Scan(&user.Name, &user.Password, &user.Email, &user.Admin)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error while trying to execute query: %s", err)
	}
	return &user, nil
}

func (repo *UserService) VerifyUserWithPassword(email string, password string) (bool, *model.User, error) {
	user, err := repo.GetByEmail(email)
	if err != nil {
		return false, nil, err
	}
	fmt.Println(user.Password)
	fmt.Println(password)
	if user.Password == password {
		return true, user, nil
	} else {
		return false, nil, nil
	}
}
