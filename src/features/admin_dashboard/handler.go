package admin_dashboard

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/db/repositories"
	"github.com/michelm117/cycling-coach-lab/features/admin_dashboard/views"
	"github.com/michelm117/cycling-coach-lab/models"
	"github.com/michelm117/cycling-coach-lab/utils"
)

type AdminDashboardHandler struct {
	repo   *repositories.UsersRepository
	repoTasks *repositories.TasksRepository
	logger *zap.SugaredLogger
}

func NewAdminDashboardHandler(
	repo *repositories.UsersRepository,
	repoTasks *repositories.TasksRepository,
	logger *zap.SugaredLogger,
) AdminDashboardHandler {
	return AdminDashboardHandler{repo: repo}
}

func (h AdminDashboardHandler) ListUsers(c echo.Context) error {
	users, err := h.repo.GetAllUsers()
	if err != nil {
		fmt.Println("error when looking for all users:" + err.Error())
	}
	return utils.Render(c, views.AdminDashboard(users))
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

	return utils.Render(c, views.AdminDashboard(users))
}

func (h AdminDashboardHandler) AddUser(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")

	newUser := models.User{
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
	return utils.Render(c, views.AdminDashboard(users))
}

func (h AdminDashboardHandler) ListTasks(c echo.Context) error {
	// return c.String(http.StatusOK, "ListTasks")
	fmt.Println(h)
	fmt.Println(h.repoTasks)
	_, err := h.repoTasks.GetAllTasks()
	if err != nil {
		fmt.Println("error when looking for all tasks:" + err.Error())
	}
	return c.String(http.StatusOK, "ListTasks")
	//return utils.Render(c, views.AdminTasks(tasks))
}

func (h AdminDashboardHandler) AddTask(c echo.Context) error {
	return c.String(http.StatusOK, "AddTask")
}

func (h AdminDashboardHandler) DeleteTask(c echo.Context) error {
	return c.String(http.StatusOK, "DeleteTask")
}