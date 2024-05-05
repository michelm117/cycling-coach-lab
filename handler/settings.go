package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/db"
	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/utils"
	"github.com/michelm117/cycling-coach-lab/views/pages"
)

type SettingsHandler struct {
	migrator db.Migrator
	logger   *zap.SugaredLogger
}

func NewSettingsHandler(
	migtator db.Migrator,
	logger *zap.SugaredLogger,
) SettingsHandler {
	return SettingsHandler{
		migrator: migtator,
		logger:   logger,
	}
}

func (h *SettingsHandler) RenderSettingsPage(c echo.Context) error {
	au := c.(model.AuthenticatedContext).User
	return Render(c, pages.SettingsPage(au, GetTheme(c)))
}

func (h *SettingsHandler) RenderSettingsView(c echo.Context) error {
	return Render(c, pages.SettingsView(GetTheme(c)))
}

func (h *SettingsHandler) Reset(c echo.Context) error {
	au := c.(model.AuthenticatedContext).User
	// TODO: Implement real authorization
	if au.Role != "admin" {
		return utils.Warning("You are not authorized to access this page")
	}

	if err := h.migrator.Reset(); err != nil {
		h.logger.Error(err)
		return utils.Danger(err.Error())
	}

	return utils.Success(c, "Anwendung erfolgreich zurückgesetzt")
}

func (h *SettingsHandler) SetTheme(c echo.Context) error {
	theme := c.FormValue("theme")

	cookie := new(http.Cookie)
	cookie.Name = "theme"
	cookie.Value = theme
	cookie.Path = "/"
	c.SetCookie(cookie)

	return utils.Success(c, "Theme erfolgreich geändert")
}
