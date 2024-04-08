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

	app.Logger.Fatal(app.Start("localhost:" + port))
}

func Setup(app *echo.Echo, db *sql.DB, logger *zap.SugaredLogger) {
	app.Use(middleware.Logger())
	if os.Getenv("ENV") == "production" {
		app.Use(middleware.Recover())
	}

	app.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/users")
	})

	utilsHandler := handler.NewUtilsHandler(db)
	app.GET("/health", utilsHandler.HealthCheck)
	app.GET("/version", utilsHandler.Version)

	userService := services.NewUserService(db, logger)
	dashboardHandler := handler.NewAdminDashboardHandler(userService, logger)

	group := app.Group("/users")
	group.GET("", dashboardHandler.ListUsers)
	group.POST("", dashboardHandler.AddUser)
	group.DELETE("/:id", dashboardHandler.DeleteUser)
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
