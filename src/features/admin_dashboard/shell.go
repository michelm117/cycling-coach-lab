package admin_dashboard

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/db/repositories"
)

func Setup(app *echo.Echo, logger *zap.SugaredLogger, repo *repositories.UsersRepository, repoTasks *repositories.TasksRepository) {
	handler := NewAdminDashboardHandler(repo, repoTasks, logger)

	group := app.Group("/users")
	group.POST("/add", handler.AddUser)
	group.GET("", handler.ListUsers)
	group.DELETE("/delete/*", handler.DeleteUser)

	groupTasks := app.Group("/tasks")
	groupTasks.POST("/add", handler.AddTask)
	groupTasks.GET("", handler.ListTasks)
	groupTasks.DELETE("/delete/*", handler.DeleteTask)
}
