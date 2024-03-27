package handler

import (
	"context"
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type HealthCheckHandler struct {
	db *sql.DB
}

func NewHealthCheckHandler(db *sql.DB) HealthCheckHandler {
	return HealthCheckHandler{
		db: db,
	}
}

func (h HealthCheckHandler) Check(c echo.Context) error {
	reqCtx := c.Request().Context()
	ctx, cancel := context.WithTimeout(reqCtx, 2*time.Second)
	defer cancel()

	err := h.db.PingContext(ctx)
	if err != nil {
		return c.String(http.StatusFailedDependency, "No conection to database")
	}

	return c.String(http.StatusOK, "Service is healthy!")
}
