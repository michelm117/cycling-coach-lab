package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/db"
	"github.com/michelm117/cycling-coach-lab/handlers"
	"github.com/michelm117/cycling-coach-lab/middlewares"
)

func main() {

	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.OpenDB()

	// Init logger
	sugar := initLogger()
	app := echo.New()

	// Middlewares
	app.Use(middlewares.RequestLogger(sugar))

	// Routes
	userHandler := handlers.UserHandler{}
	handlers.SetupRoutes(app, &userHandler)

	// Serve static files
	app.Static("/assets", "assets")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app.Logger.Fatal(app.Start(":" + "8080"))
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
