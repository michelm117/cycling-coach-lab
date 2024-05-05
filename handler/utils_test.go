package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/michelm117/cycling-coach-lab/handler"
)

func TestGetTheme(t *testing.T) {
	t.Run("Theme not set", func(t *testing.T) {
		// Create a request
		req := httptest.NewRequest(http.MethodGet, "/settings/view", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		assert.Equal(t, "dark", handler.GetTheme(c))
	})

	t.Run("Theme set", func(t *testing.T) {
		// Create a request
		req := httptest.NewRequest(http.MethodGet, "/settings/view", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("theme", "light")

		assert.Equal(t, "light", handler.GetTheme(c))
	})
}
