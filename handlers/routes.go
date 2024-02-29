package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(app *echo.Echo, h *UserHandler) {
	healthcheckRoute(app)
	userRoutes(app, h)
}

func healthcheckRoute(app *echo.Echo) {
	app.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service is healthy!")
	})
}

func userRoutes(app *echo.Echo, h *UserHandler) {
	group := app.Group("/users")

	group.POST("/add", h.HandlerAddUser)
	group.GET("", h.HandlerShowUsers)
	group.GET("/details/:id", h.HandlerShowUserById)
	//TODO: this is a post because I couldnt get the values with a delete because we dont have a unique key rn
	group.DELETE("/delete/*", h.HandleDeleteUser)
}
