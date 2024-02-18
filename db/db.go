package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/stdlib"
)

func ConnectToDatabase() *sql.DB {
	psqlInfo, err := buildPsqlInfo()
	if err != nil {
		log.Fatal("Error:", err)
	}
	println(psqlInfo)
	db, err := sql.Open("pgx", psqlInfo)
	println(db)
	if err != nil {
		log.Fatal("Error while connecting to db cause: " + err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Error while pinging to db cause: " + err.Error())
	}

	return db
}

func buildPsqlInfo() (string, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Check if required environment variables are set
	if host == "" {
		return "", fmt.Errorf("DB_HOST environment variable is required")
	}
	if port == "" {
		return "", fmt.Errorf("DB_PORT environment variable is required")
	}
	if user == "" {
		return "", fmt.Errorf("DB_USER environment variable is required")
	}
	if password == "" {
		return "", fmt.Errorf("DB_PASSWORD environment variable is required")
	}
	if dbname == "" {
		return "", fmt.Errorf("DB_NAME environment variable is required")
	}

	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)

	return psqlInfo, nil
}
