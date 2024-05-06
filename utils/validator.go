package utils

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/michelm117/cycling-coach-lab/model"
)

type Validator interface {
	ValidateNonEmptyStringField(fieldName, value string) error
	ValidateRole(role string) error
	ValidateEmail(email string) error
	ValidatePassword(password, confirmPassword string) error
	CreateValidUser(firstname, lastname, role, email, dateOfBirthStr, password, confirmPassword string) (*model.User, error)
}

func NewValidator(cryptoer Cryptoer) Validator {
	return &validator{cryptoer: cryptoer}
}

type validator struct {
	cryptoer Cryptoer
}

func (v *validator) ValidateNonEmptyStringField(fieldName, value string) error {
	if value == "" {
		return fmt.Errorf("Invalid %s", fieldName)
	}
	return nil
}

func (v *validator) ValidateRole(role string) error {
	if role != "admin" && role != "user" {
		return fmt.Errorf("Invalid role")
	}
	return nil
}

func (v *validator) ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("Invalid email")
	}
	return nil
}

func (v *validator) ValidatePassword(password, confirmPassword string) error {
	if password != confirmPassword {
		return fmt.Errorf("Passwords do not match")
	}
	return nil
}

func (v *validator) CreateValidUser(firstname, lastname, role, email, dateOfBirthStr, password, confirmPassword string) (*model.User, error) {
	if err := v.ValidateNonEmptyStringField("first name", firstname); err != nil {
		return nil, err
	}

	if err := v.ValidateNonEmptyStringField("last name", lastname); err != nil {
		return nil, err
	}

	if err := v.ValidateRole(role); err != nil {
		return nil, err
	}

	if err := v.ValidateEmail(email); err != nil {
		return nil, err
	}

	dateOfBirth, err := time.Parse("2006-01-02", dateOfBirthStr)
	if err != nil {
		return nil, err
	}

	if err := v.ValidatePassword(password, confirmPassword); err != nil {
		return nil, err
	}

	hashedPassword, err := v.cryptoer.GenerateFromPassword([]byte(password))
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Firstname:    firstname,
		Lastname:     lastname,
		Email:        email,
		DateOfBirth:  dateOfBirth,
		Role:         role,
		Status:       "active",
		PasswordHash: string(hashedPassword),
	}

	return user, nil
}
