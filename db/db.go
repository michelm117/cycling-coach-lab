package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/stdlib"
)

func ConnectToDatabase() *sql.DB {

	env, err := GetDatabaseEnv()
	if err != nil {
		log.Fatal("Error:", err)
	}

	db, err := sql.Open("pgx", env.Address)
	if err != nil {
		log.Fatal("Error while connecting to db cause: " + err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error while pinging to db cause: %s", err.Error())
	}

	// ToDo: migration
	return db
}
