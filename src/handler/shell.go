package handler

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/db/repositories"
)

func Setup(app *echo.Echo, logger *zap.SugaredLogger, repo *repositories.UsersRepository) {
	handler := NewAdminDashboardHandler(repo, logger)

	group := app.Group("/users")
	group.POST("/add", handler.AddUser)
	group.GET("", handler.ListUsers)
	group.DELETE("/delete/*", handler.DeleteUser)
}
