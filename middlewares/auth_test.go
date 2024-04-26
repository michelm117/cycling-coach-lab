package middlewares_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/michelm117/cycling-coach-lab/middlewares"
	"github.com/michelm117/cycling-coach-lab/mocks"
	"github.com/michelm117/cycling-coach-lab/model"
)

func setupEchoContext() (e *echo.Echo, c echo.Context, rec *httptest.ResponseRecorder) {
	e = echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	return
}

func createValidSession() *sessions.Session {
	return &sessions.Session{
		Values: make(map[interface{}]interface{}),
		Options: &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		},
		IsNew: false,
		ID:    "valid_session_i",
	}
}

func TestAuthenticationMiddleware_ValidSession(t *testing.T) {
	// Setup
	_, c, rec := setupEchoContext()

	// Mock session service
	session := createValidSession()
	session.Values["sessionId"] = "valid_session_id"

	ctrl := gomock.NewController(t)

	mb := mocks.NewMockBrowserSessionManager(ctrl)
	mb.EXPECT().Get(gomock.Eq(c)).Return(session, nil)

	ms := mocks.NewMockSessionServicer(ctrl)
	ms.EXPECT().AuthenticateUserBySessionID(gomock.Eq("valid_session_id")).Return(&model.User{ID: 1}, nil)

	// Test middleware
	authMiddleware := middlewares.Authentication(ms, mb)
	handler := authMiddleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "Authorized")
	})
	err := handler(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Authorized", rec.Body.String())
}

func TestAuthenticationMiddleware_InvalidSession(t *testing.T) {
	// Setup
	_, c, rec := setupEchoContext()

	// Mock session service without setting sessionID
	session := createValidSession()

	ctrl := gomock.NewController(t)
	mb := mocks.NewMockBrowserSessionManager(ctrl)
	mb.EXPECT().Get(gomock.Eq(c)).Return(session, nil)

	ms := mocks.NewMockSessionServicer(ctrl)
	ms.EXPECT().AuthenticateUserBySessionID(gomock.Eq("")).Return(nil, fmt.Errorf("sessionId is empty"))

	// Test middleware
	authMiddleware := middlewares.Authentication(ms, mb)
	handler := authMiddleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "Authorized")
	})
	err := handler(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/auth/login", rec.Header().Get("Location"))
	assert.Equal(t, "", rec.Body.String())
}

func TestAuthenticationMiddleware_BrowserSessionError(t *testing.T) {
	// Setup
	_, c, rec := setupEchoContext()

	// Mock session service
	session := createValidSession()

	ctrl := gomock.NewController(t)
	mb := mocks.NewMockBrowserSessionManager(ctrl)
	mb.EXPECT().Get(gomock.Eq(c)).Return(session, assert.AnError)

	// Test middleware
	authMiddleware := middlewares.Authentication(nil, mb)
	handler := authMiddleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "Authorized")
	})
	err := handler(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/auth/login", rec.Header().Get("Location"))
	assert.Equal(t, "", rec.Body.String())
}
