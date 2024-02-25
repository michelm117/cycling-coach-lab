package db

import (
	"fmt"
	"os"
)

type DatabaseEnv struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Address  string
}

func GetDatabaseEnv() (*DatabaseEnv, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Check if required environment variables are set
	if host == "" {
		return nil, fmt.Errorf("DB_HOST environment variable is required")
	}
	if port == "" {
		return nil, fmt.Errorf("DB_PORT environment variable is required")
	}
	if user == "" {
		return nil, fmt.Errorf("DB_USER environment variable is required")
	}
	if password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}
	if dbname == "" {
		return nil, fmt.Errorf("DB_NAME environment variable is required")
	}

	adress := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	return &DatabaseEnv{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     dbname,
		Address:  adress,
	}, nil
}
