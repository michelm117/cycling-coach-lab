package model

import (
	"time"
)

// User struct represents the users table in the database
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	PasswordHash string    `json:"password_hash"`
	Status       string    `json:"status"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
