package handler

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type UtilsHandler struct {
	db *sql.DB
}

func NewUtilsHandler(db *sql.DB) UtilsHandler {
	return UtilsHandler{
		db: db,
	}
}

func (h UtilsHandler) HealthCheck(c echo.Context) error {
	reqCtx := c.Request().Context()
	ctx, cancel := context.WithTimeout(reqCtx, 2*time.Second)
	defer cancel()

	err := h.db.PingContext(ctx)
	if err != nil {
		return c.String(http.StatusFailedDependency, "No conection to database")
	}

	return c.String(http.StatusOK, "Service is healthy!")
}

func (h UtilsHandler) Version(c echo.Context) error {
	version := os.Getenv("VERSION")
	return c.String(http.StatusOK, version)
}
