package handlers

import "github.com/labstack/echo/v4"

func SetupRoutes(app *echo.Echo, h *UserHandler) {
	group := app.Group("/users")

	group.POST("/add", h.HandlerAddUser)
	group.GET("", h.HandlerShowUsers)
	group.GET("/details/:id", h.HandlerShowUserById)
}
