package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/views/admin_dashboard"
	"github.com/michelm117/cycling-coach-lab/views/utils"
)

type AdminDashboardHandler struct {
	repo   *services.UserService
	logger *zap.SugaredLogger
}

func NewAdminDashboardHandler(
	repo *services.UserService,
	logger *zap.SugaredLogger,
) AdminDashboardHandler {
	return AdminDashboardHandler{repo: repo, logger: logger}
}

func (h AdminDashboardHandler) ListUsers(c echo.Context) error {
	users, err := h.repo.GetAllUsers()
	if err != nil {
		fmt.Println("error when looking for all users:" + err.Error())
	}
	return Render(c, admin_dashboard.Index(users), http.StatusOK)
}

func (h AdminDashboardHandler) DeleteUser(c echo.Context) error {
	encodedEmail := c.ParamValues()
	email, err := url.QueryUnescape(string(encodedEmail[0]))
	if err != nil {
		msg := fmt.Sprintf("Could not decode url encoded email: '%s'", email)
		return Render(c, utils.AlertError(msg), http.StatusBadRequest)
	}

	if err := h.repo.DeleteUser(email); err != nil {
		msg := fmt.Sprintf("Could not delete user with email '%s'", email)
		return Render(c, utils.AlertError(msg), http.StatusBadRequest)
	}

	users, err := h.repo.GetAllUsers()
	if err != nil {
		return Render(c, utils.AlertError(err.Error()), http.StatusBadRequest)
	}

	return Render(c, admin_dashboard.UserTable(users), http.StatusOK)
}

func (h AdminDashboardHandler) AddUser(c echo.Context) error {
	firstname := c.FormValue("firstname")
	lastname := c.FormValue("lastname")
	email := c.FormValue("email")
	role := c.FormValue("role")

	dateOfBirthStr := c.FormValue("dateOfBirth")
	dateOfBirth, err := time.Parse("2006-01-02", dateOfBirthStr)

	password := c.FormValue("password")
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	if err != nil {
		return err
	}

	newUser := model.User{
		Firstname:    firstname,
		Lastname:     lastname,
		DateOfBirth:  dateOfBirth,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		Status:       "active",
	}

	_, err = h.repo.AddUser(newUser)
	if err != nil {
		h.logger.Warnf("Error while adding user: %s", err.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	users, _ := h.repo.GetAllUsers()
	return Render(c, admin_dashboard.UserTable(users), http.StatusOK)
}
