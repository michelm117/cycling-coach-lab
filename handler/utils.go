package handler

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, component templ.Component, resStatus int) error {
	c.Response().Status = resStatus // fix: maybe there is a better way to set the statuscode
	return component.Render(c.Request().Context(), c.Response())
}
