package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/michelm117/cycling-coach-lab/models"
)

var DB *sql.DB

func OpenDB() error {
	psqlInfo := "postgres://postgres:postgres@localhost:5432/postgres"
	var err error
	DB, err = sql.Open("pgx", psqlInfo)
	println("what is the db object:")
	println(DB)
	if err != nil {
		fmt.Printf("Error while connecting to db cause: " + err.Error())
		return err
	}
	if err := DB.Ping(); err != nil {
		fmt.Printf("Error while pinning to db cause: " + err.Error())

		return err
	}

	return nil
}

func GetUserByName(name string) (models.User, error) {
	var user models.User

	row := DB.QueryRow("SELECT username, email FROM users WHERE users.username = $1", name)
	err := row.Scan(&user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found")
		}
		log.Println("Error when trying to execute query")
		log.Println(err)
		return models.User{}, err
	}

	return user, nil
}

func AddUser(user models.User) error {
	println("we are fucking adding a user")

	_, err := DB.Exec("INSERT INTO users (username, email) VALUES ($1, $2)", user.Name, user.Email)
	if err != nil {
		log.Println("Error when trying to execute query")
		log.Println(err)
		return err
	}

	println("we fucking added a user")
	return nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	println(&DB)
	rows, err := DB.Query("SELECT username, email FROM users")
	if err != nil {
		log.Println("Error when trying to execute query")
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Name, &user.Email)
		if err != nil {
			log.Println("Error scanning row")
			log.Println(err)
			continue
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error after iterating over all rows")
		log.Println(err)
		return nil, err
	}

	defer rows.Close()
	return users, nil
}
