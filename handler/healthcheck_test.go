package handler_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/michelm117/cycling-coach-lab/handler"
)

func TestHealthCheckHandler(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	handler := handler.NewUtilsHandler(db)

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	assert.NoError(t, handler.HealthCheck(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestVersionHandler(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Set environment variable VERSION for testing
	os.Setenv("VERSION", "1.0.0")
	defer os.Unsetenv("VERSION")

	handler := handler.NewUtilsHandler(db)

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	assert.NoError(t, handler.Version(c))
	assert.Equal(t, "1.0.0", rec.Body.String())
}
