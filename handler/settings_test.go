package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"

	"github.com/michelm117/cycling-coach-lab/handler"
	"github.com/michelm117/cycling-coach-lab/mocks"
	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/test_utils"
)

func TestRenderSettingsPage(t *testing.T) {
	au := model.User{ID: 1, Firstname: "John", Lastname: "Doe", Email: "john@doe.com"}
	handler := handler.NewSettingsHandler(nil, nil, nil)

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/settings", nil)
	rec := httptest.NewRecorder()
	c := model.AuthenticatedContext{Context: echo.New().NewContext(req, rec), User: &au}

	// Call the handler
	assert.NoError(t, handler.RenderSettingsPage(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	test_utils.MakeSnapshot(t, rec.Body.String())
}

func TestRenderSettingsView(t *testing.T) {
	handler := handler.NewSettingsHandler(nil, nil, nil)

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/settings/view", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Call the handler
	assert.NoError(t, handler.RenderSettingsView(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	test_utils.MakeSnapshot(t, rec.Body.String())
}

func TestReset(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mm := mocks.NewMockMigrator(ctrl)
	es := mocks.NewMockEmailServicer(ctrl)

	t.Run("Not authorized", func(t *testing.T) {
		au := model.User{ID: 1, Firstname: "John", Lastname: "Doe", Email: "john@doe.com", Role: "athlete"}
		handler := handler.NewSettingsHandler(nil, nil, nil)

		// Create a request
		req := httptest.NewRequest(http.MethodGet, "/settings/reset", nil)
		rec := httptest.NewRecorder()
		c := model.AuthenticatedContext{Context: echo.New().NewContext(req, rec), User: &au}

		// Call the handler
		assert.ErrorContains(t, handler.Reset(c), "You are not authorized to access this page")
		assert.Equal(t, http.StatusOK, rec.Code)
		test_utils.MakeSnapshot(t, rec.Body.String())
	})

	t.Run("migrator error", func(t *testing.T) {
		au := model.User{ID: 1, Firstname: "John", Lastname: "Doe", Email: "john@doe.com", Role: "admin"}

		mm.EXPECT().Reset().Return(assert.AnError)

		handler := handler.NewSettingsHandler(es, mm, logger)

		// Create a request
		req := httptest.NewRequest(http.MethodGet, "/settings/reset", nil)
		rec := httptest.NewRecorder()
		c := model.AuthenticatedContext{Context: echo.New().NewContext(req, rec), User: &au}

		// Call the handler
		assert.Error(t, handler.Reset(c))
		assert.Equal(t, http.StatusOK, rec.Code)
		test_utils.MakeSnapshot(t, rec.Body.String())
	})

	t.Run("Success", func(t *testing.T) {
		au := model.User{ID: 1, Firstname: "John", Lastname: "Doe", Email: "john@doe.com", Role: "admin"}

		mm.EXPECT().Reset().Return(nil)

		handler := handler.NewSettingsHandler(es, mm, logger)

		// Create a request
		req := httptest.NewRequest(http.MethodGet, "/settings/reset", nil)
		rec := httptest.NewRecorder()
		c := model.AuthenticatedContext{Context: echo.New().NewContext(req, rec), User: &au}

		// Call the handler
		assert.NoError(t, handler.Reset(c))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "/auth/login", rec.Header().Get("HX-Redirect"))
		test_utils.MakeSnapshot(t, rec.Body.String())
	})
}

func TestSetTheme(t *testing.T) {
	handler := handler.NewSettingsHandler(nil, nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("theme=aqua"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()

	c := model.AuthenticatedContext{Context: echo.New().NewContext(req, rec)}

	err := handler.SetTheme(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	cookie := rec.Result().Cookies()[0]
	assert.Equal(t, "theme", cookie.Name)
	assert.Equal(t, "aqua", cookie.Value)
	assert.Equal(t, "/", cookie.Path)
	test_utils.MakeSnapshot(t, rec.Body.String())
}
