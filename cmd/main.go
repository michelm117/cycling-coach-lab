package main

import (
	"database/sql"
	"fmt"
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

	address := fmt.Sprintf(":%v", port)
	if os.Getenv("ENV") == "development" {
		address = fmt.Sprintf("localhost:%v", port)
	}
	logger.Infof("Starting server on %v", address)
	app.Logger.Fatal(app.Start(address))
}

func Setup(app *echo.Echo, db *sql.DB, logger *zap.SugaredLogger) {
	app.Use(middleware.Logger())
	if os.Getenv("ENV") == "production" {
		app.Use(middleware.Recover())
	}
	app.HTTPErrorHandler = customErrorHandler

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

func customErrorHandler(err error, c echo.Context) {
	// Attempt casting the error as a Toast.
	te, ok := err.(handler.Toast)

	// If it canot be cast as a Toast, it must be some other error
	// we did not handle. We will handle it here and return a more
	// generic error message. We don't want system errors to bleed
	// through to the user.
	if !ok {
		fmt.Println(err)
		te = handler.Danger("there has been an unexpected error")
	}

	// If not a success error (weird right) set the HX-Swap header to `none`.
	if te.Level != handler.SUCCESS {
		c.Response().Header().Set("HX-Reswap", "none")
	}

	// Set the HX-Trigger header
	te.SetHXTriggerHeader(c)
}
