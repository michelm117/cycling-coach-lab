package handler

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

func GetTheme(c echo.Context) string {
	theme, ok := c.Get("theme").(string)
	if !ok {
		return "dark-light"
	}
	return theme
}
