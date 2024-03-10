package shell

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/db/repositories"
	"github.com/michelm117/cycling-coach-lab/features/admin_dashboard"
	"github.com/michelm117/cycling-coach-lab/shell/middlewares"
)

func Setup(app *echo.Echo, db *sql.DB, logger *zap.SugaredLogger) {
	app.Use(middlewares.RequestLogger(logger))

	app.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/users")
	})

	app.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service is healthy!")
	})

	usersRepository := repositories.NewUsersRepository(db, logger)

	admin_dashboard.Setup(app, logger, usersRepository)
}
