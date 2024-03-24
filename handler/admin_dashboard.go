package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/views/admin_dashboard"
)

type AdminDashboardHandler struct {
	repo   *services.UserService
	logger *zap.SugaredLogger
}

func NewAdminDashboardHandler(
	repo *services.UserService,
	logger *zap.SugaredLogger,
) AdminDashboardHandler {
	return AdminDashboardHandler{repo: repo}
}

func (h AdminDashboardHandler) ListUsers(c echo.Context) error {
	users, err := h.repo.GetAllUsers()
	if err != nil {
		fmt.Println("error when looking for all users:" + err.Error())
	}
	return Render(c, admin_dashboard.Index(users))
}

func (h AdminDashboardHandler) DeleteUser(c echo.Context) error {
	email := c.ParamValues()
	emailOfUser := strings.Replace(email[0], "%40", "@", -1)
	userToBeDeleted, err := h.repo.GetByEmail(emailOfUser)
	if err != nil {
		return err
	}
	h.repo.DeleteUser(*userToBeDeleted)
	users, _ := h.repo.GetAllUsers()
	for _, t := range users {
		println(t.Email)
		println(t.Name)
	}

	return Render(c, admin_dashboard.UserTable(users))
}

func (h AdminDashboardHandler) AddUser(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")

	newUser := model.User{
		Name:  name,
		Email: email,
	}
	_, err := h.repo.AddUser(newUser)
	if err != nil {
		h.logger.Warnf("Error while adding user: %s", err.Error())
		// if strings.Contains(err.Error(), "duplicate") {
		//   return render(c, components.EmailTaken(newUser))
		// }
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	users, _ := h.repo.GetAllUsers()
	return Render(c, admin_dashboard.UserTable(users))
}
