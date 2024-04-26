package handler

import (
	"fmt"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/utils"
	"github.com/michelm117/cycling-coach-lab/views/pages"
)

type UserManagementHandler struct {
	userServicer services.UserServicer
	cryptoer     utils.Cryptoer
	logger       *zap.SugaredLogger
}

func NewUserManagementHandler(
	userServicer services.UserServicer,
	cryptoer utils.Cryptoer,
	logger *zap.SugaredLogger,
) UserManagementHandler {
	return UserManagementHandler{userServicer: userServicer, cryptoer: cryptoer, logger: logger}
}

func (h UserManagementHandler) RenderUserTable(c echo.Context) error {
	au := c.(model.AuthenticatedContext).User

	users, err := h.userServicer.GetAllUsers()
	if err != nil {
		return utils.Warning("Could not retrieve users")
	}
	return Render(c, pages.UserManagementIndex(au, users))
}

func (h UserManagementHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.Warning("Invalid user id")
	}

	if err := h.userServicer.DeleteUser(id); err != nil {
		return utils.Warning(fmt.Sprintf("Could not delete user with id '%v'", id))
	}

	return utils.Success(c, "User deleted successfully")
}

func (h UserManagementHandler) RenderAddUser(c echo.Context) error {
	// Validate first name
	firstname := c.FormValue("firstname")
	if err := validateNonEmptyStringField("first name", firstname); err != nil {
		return err
	}

	// Validate last name
	lastname := c.FormValue("lastname")
	if err := validateNonEmptyStringField("last name", lastname); err != nil {
		return err
	}

	// Validate role
	role := c.FormValue("role")
	if err := validateRole(role); err != nil {
		return err
	}

	// Validate email
	email := c.FormValue("email")
	if err := validateEmail(email); err != nil {
		return err
	}

	// Validate date of birth
	dateOfBirthStr := c.FormValue("dateOfBirth")
	dateOfBirth, err := time.Parse("2006-01-02", dateOfBirthStr)
	if err != nil {
		return utils.Warning("Invalid date of birth")
	}

	// Validate password
	password := c.FormValue("password")
	hashedPassword, err := h.cryptoer.GenerateFromPassword([]byte(password))
	if err != nil {
		return utils.Warning("Could not add user")
	}

	// Create user object
	user := model.User{
		Firstname:    firstname,
		Lastname:     lastname,
		DateOfBirth:  dateOfBirth,
		Email:        strings.ToLower(email),
		PasswordHash: string(hashedPassword),
		Role:         role,
		Status:       "active",
	}

	// Add user
	if _, err := h.userServicer.AddUser(user); err != nil {
		return utils.Warning("Could not add user")
	}

	// Success response
	utils.Success(c, fmt.Sprintf("User '%s' added successfully", user.Email))
	return Render(c, pages.AddUserResponse(&user))
}

// Validation functions
func validateNonEmptyStringField(fieldName, value string) error {
	if value == "" {
		return utils.Warning("Invalid " + fieldName)
	}
	return nil
}

func validateRole(role string) error {
	if role != "admin" && role != "user" {
		return utils.Warning("Invalid role")
	}
	return nil
}

func validateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return utils.Warning("Invalid email")
	}
	return nil
}
