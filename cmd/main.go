package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/michelm117/cycling-coach-lab/handlers"
	"github.com/michelm117/cycling-coach-lab/middlewares"
	"go.uber.org/zap"
)

func main() {
	// Init logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

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
	app.Logger.Fatal(app.Start(":" + port))
}
