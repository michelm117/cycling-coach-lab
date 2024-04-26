package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/utils"
)

func Authentication(sessionService services.SessionServicer, browserSessionManager utils.BrowserSessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			browserSession, err := browserSessionManager.Get(c)
			if err != nil {
				return c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
			}

			sessionID, _ := browserSession.Values["sessionId"].(string)
			user, err := sessionService.AuthenticateUserBySessionID(sessionID)
			if err != nil {
				return c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
			}

			// Extend the echo.Context with the authenticated user
			authenticatedContext := model.AuthenticatedContext{
				User:    user,
				Context: c,
			}
			return next(authenticatedContext)
		}
	}
}
