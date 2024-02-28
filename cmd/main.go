package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/db"
	"github.com/michelm117/cycling-coach-lab/handlers"
	"github.com/michelm117/cycling-coach-lab/middlewares"
	"github.com/michelm117/cycling-coach-lab/repositories"
)

func main() {

	logger := initLogger()

	db := db.ConnectToDatabase(logger)

	app := echo.New()
	app.Use(middlewares.RequestLogger(logger))
	userRepository := repositories.NewUserRepository(db, logger)
	userHandler := handlers.NewUserHandler(userRepository)
	handlers.SetupRoutes(app, &userHandler)

	// Serve static files
	app.Static("/assets", "assets")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Logger.Fatal(app.Start(":" + port))
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
