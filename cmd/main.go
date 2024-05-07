package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/db"
	"github.com/michelm117/cycling-coach-lab/handler"
	"github.com/michelm117/cycling-coach-lab/middlewares"
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
	database := db.ConnectToDatabase(logger)
	migrator := db.NewMigrator(database, "migrations", logger)
	if err := migrator.Up(); err != nil {
		log.Fatal(err)
	}
	app := echo.New()

	// Serve static files
	assetsPath := path.Join(utils.GetProjectRoot(), "assets")
	logger.Infof("Serving static files from: %s", assetsPath)
	app.Static("/assets", assetsPath)

	Setup(app, database, migrator, logger)
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

func Setup(app *echo.Echo, db *sql.DB, migrator db.Migrator, logger *zap.SugaredLogger) {
	if os.Getenv("ENV") == "production" {
		app.Use(middleware.Logger())
		app.Use(middleware.Recover())
	}

	app.HTTPErrorHandler = customErrorHandler

	secret := os.Getenv("SESSION_SECRET")
	app.Use(session.Middleware(sessions.NewCookieStore([]byte(secret))))

	// middleware that check if theme cookie is set and set the theme in the context
	// TODO: move to middlewares/theme.go
	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			themeCookie, err := c.Cookie("theme")
			if err == nil {
				c.Set("theme", themeCookie.Value)
			} else {
				c.Set("theme", "dark")
			}
			return next(c)
		}
	})

	app.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/users")
	})

	// Initialize Casbin Enforcer
	enforcer := utils.NewCasbinEnforcer("./casbin/casbin_auth_model.conf", "./casbin/casbin_auth_policy.csv")
	// Health check and version endpoints
	utilsHandler := handler.NewUtilsHandler(db)
	app.GET("/health", utilsHandler.HealthCheck)
	app.GET("/version", utilsHandler.Version)

	cryptoer := utils.NewCrypto()
	validator := utils.NewValidator(cryptoer)
	browserSessionManager := utils.NewBrowserSessionManager()
	globalSettingsServicer := services.NewGlobalSettingService(db, logger)
	userServicer := services.NewUserServicer(db, logger)
	sessionService := services.NewSessionServicer(db, logger)
	sessionService.ScheduleSessionCleanup()

	setupHandler := handler.NewSetupHandler(globalSettingsServicer, userServicer, validator, logger)
	app.GET("/setup", setupHandler.RenderSetup)
	app.POST("/setup", setupHandler.Setup)

	authHandler := handler.NewAuthHandler(userServicer, sessionService, globalSettingsServicer, browserSessionManager, cryptoer, logger)
	authRoute := app.Group("/auth")
	authRoute.GET("/login", authHandler.RenderLogin)
	authRoute.POST("/login", authHandler.Login)
	authRoute.POST("/logout", authHandler.Logout)

	userManagementHandler := handler.NewUserManagementHandler(userServicer, validator, logger)
	usersRoute := app.Group("/users")
	usersRoute.Use(middlewares.Authentication(sessionService, browserSessionManager))
	usersRoute.Use(middlewares.Autheratziation(enforcer))
	usersRoute.GET("", userManagementHandler.RenderUserManagementPage)
	usersRoute.GET("/view", userManagementHandler.RenderUserManagementView)
	usersRoute.POST("", userManagementHandler.RenderAddUser)
	usersRoute.DELETE("/:id", userManagementHandler.DeleteUser)

	settingsHandler := handler.NewSettingsHandler(migrator, logger)
	settingsRoute := app.Group("/settings")
	settingsRoute.Use(middlewares.Authentication(sessionService, browserSessionManager))
	settingsRoute.Use(middlewares.Autheratziation(enforcer))
	settingsRoute.GET("", settingsHandler.RenderSettingsPage)
	settingsRoute.GET("/view", settingsHandler.RenderSettingsView)
	settingsRoute.POST("/reset", settingsHandler.Reset)
	settingsRoute.POST("/theme", settingsHandler.SetTheme)
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

// TODO: move to utils/toast.go
func customErrorHandler(err error, c echo.Context) {
	// Attempt casting the error as a Toast.
	te, ok := err.(utils.Toast)

	// If it canot be cast as a Toast, it must be some other error
	// we did not handle. We will handle it here and return a more
	// generic error message. We don't want system errors to bleed
	// through to the user.
	if !ok {
		te = utils.Danger("there has been an unexpected error")
		fmt.Println("Unexpected error:", err.Error())
	}

	// If not a success error (weird right) set the HX-Swap header to `none`.
	if te.Level != utils.SUCCESS {
		c.Response().Header().Set("HX-Reswap", "none")
	}

	// Set the HX-Trigger header
	te.SetHXTriggerHeader(c)
}
