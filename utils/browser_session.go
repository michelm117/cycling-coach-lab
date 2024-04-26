package utils

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type BrowserSessionManager interface {
	Get(c echo.Context) (*sessions.Session, error)
	SaveSession(c echo.Context, sessionID string) error
	DeleteSession(c echo.Context) (string, error)
}

type BrowserSessionManage struct{}

const (
	sessionMaxAgeSeconds = 86400 * 7
)

func NewBrowserSessionManager() BrowserSessionManager {
	return &BrowserSessionManage{}
}

func (m *BrowserSessionManage) Get(c echo.Context) (*sessions.Session, error) {
	return session.Get("cyclinet-coach-lab", c)
}

func (m *BrowserSessionManage) SaveSession(c echo.Context, sessionID string) error {
	browserSession, err := session.Get("cyclinet-coach-lab", c)
	if err != nil {
		return err
	}

	// Configure session options
	browserSession.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionMaxAgeSeconds,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") != "development",
	}
	browserSession.Values["sessionId"] = sessionID

	if err := browserSession.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return nil
}

func (m *BrowserSessionManage) DeleteSession(c echo.Context) (string, error) {
	browserSession, err := m.Get(c)
	if err != nil {
		return "", err
	}
	sessionID := browserSession.Values["sessionId"].(string)

	browserSession.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") != "development",
	}

	if err := browserSession.Save(c.Request(), c.Response()); err != nil {
		return "", err
	}

	return sessionID, nil
}
