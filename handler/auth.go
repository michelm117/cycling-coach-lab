package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/utils"
	"github.com/michelm117/cycling-coach-lab/views/pages"
)

const (
	sessionCookieName = "cycling-coach-lab"
)

type AuthHandler struct {
	userServicer          services.UserServicer
	sessionServicer       services.SessionServicer
	globalSettingServicer services.GlobalSettingServicer
	browserSessionManager utils.BrowserSessionManager
	crptoer               utils.Cryptoer
	logger                *zap.SugaredLogger
}

func NewAuthHandler(
	userServicer services.UserServicer,
	sessionServicer services.SessionServicer,
	globalSettingServicer services.GlobalSettingServicer,
	browserSessionManager utils.BrowserSessionManager,
	crptoer utils.Cryptoer,
	logger *zap.SugaredLogger,
) AuthHandler {
	return AuthHandler{
		userServicer:          userServicer,
		sessionServicer:       sessionServicer,
		globalSettingServicer: globalSettingServicer,
		browserSessionManager: browserSessionManager,
		crptoer:               crptoer,
		logger:                logger,
	}
}

func (h AuthHandler) RenderLogin(c echo.Context) error {
	if !h.globalSettingServicer.IsAppInitialized() {
		return c.Redirect(http.StatusTemporaryRedirect, "/setup")
	}

	return Render(c, pages.Login())
}

func (h AuthHandler) Login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := h.userServicer.GetByEmail(email)
	if err != nil {
		h.logger.Error(err)
		return utils.Warning("Invalid credentials")
	}

	if user.Status != "active" {
		return utils.Warning("Invalid credentials")
	}

	if err := h.crptoer.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return utils.Warning("Invalid credentials")
	}

	// Save session ID in the database and in the browser session
	sessionID, err := h.sessionServicer.SaveSession(user.ID)
	if err != nil {
		h.logger.Error(err)
		return utils.Warning("Invalid credentials")
	}

	if err = h.browserSessionManager.SaveSession(c, sessionID); err != nil {
		h.logger.Error(err)
		return utils.Warning("Invalid credentials")
	}

	c.Response().Header().Add("HX-Redirect", "/users")
	return nil
}

func (h AuthHandler) Logout(c echo.Context) error {
	sessionID, err := h.browserSessionManager.DeleteSession(c)
	if err != nil {
		return utils.Warning("Internal server error")
	}

	if err := h.sessionServicer.DeleteSession(sessionID); err != nil {
		return utils.Warning("Internal server error")
	}

	c.Response().Header().Add("HX-Redirect", "/auth/login")
	return nil
}
