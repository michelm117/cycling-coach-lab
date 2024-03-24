package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/db"
	"github.com/michelm117/cycling-coach-lab/handler"
	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/utils"
)

func main() {
	err := utils.CheckForRequiredEnvVars()
	if err != nil {
		log.Fatal("Error:", err)
	}

	logger := initLogger()
	logger.Infof("Starting server in `%s` mode", os.Getenv("ENV"))
	db := db.ConnectToDatabase(logger)
	app := echo.New()

	// Serve static files
	assetsPath := path.Join(utils.GetProjectRoot(), "assets")
	logger.Infof("Serving static files from: %s", assetsPath)
	app.Static("/assets", assetsPath)

	Setup(app, db, logger)
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Logger.Fatal(app.Start(":" + port))
}

func Setup(app *echo.Echo, db *sql.DB, logger *zap.SugaredLogger) {
	app.Use(middleware.Logger())

	app.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/users")
	})

	app.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service is healthy!")
	})

	userService := services.NewUserService(db, logger)
	handler := handler.NewAdminDashboardHandler(userService, logger)

	group := app.Group("/users")
	group.POST("/add", handler.AddUser)
	group.GET("", handler.ListUsers)
	group.DELETE("/delete/*", handler.DeleteUser)
}

func initLogger() *zap.SugaredLogger {
	var logger *zap.Logger
	if os.Getenv("ENV") == "development" {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()
	return logger.Sugar()
}
