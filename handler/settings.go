package handler

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/db"
	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/utils"
	"github.com/michelm117/cycling-coach-lab/views/pages"
)

type SettingsHandler struct {
	emailServicer services.EmailServicer
	migrator      db.Migrator
	logger        *zap.SugaredLogger
}

func NewSettingsHandler(
	emailServicer services.EmailServicer,
	migtator db.Migrator,
	logger *zap.SugaredLogger,
) SettingsHandler {
	return SettingsHandler{
		emailServicer: emailServicer,
		migrator:      migtator,
		logger:        logger,
	}
}

func (h *SettingsHandler) RenderSettingsPage(c echo.Context) error {
	au := c.(model.AuthenticatedContext).User
	emailSettings, err := h.emailServicer.GetEmailSettings()
	emailSettings.Password = ""

	if err == sql.ErrNoRows {
		emailSettings = &model.EmailSettings{}
		return Render(c, pages.SettingsPage(au, GetTheme(c), *emailSettings))
	}

	if err != nil {
		return utils.Danger(err.Error())
	}

	return Render(c, pages.SettingsPage(au, GetTheme(c), *emailSettings))
}

func (h *SettingsHandler) RenderSettingsView(c echo.Context) error {
	emailSettings, err := h.emailServicer.GetEmailSettings()
	emailSettings.Password = ""
	if err != nil {
		return utils.Danger(err.Error())
	}
	return Render(c, pages.SettingsView(GetTheme(c), *emailSettings))
}

func (h *SettingsHandler) SaveEmailSettings(c echo.Context) error {
	settings := model.EmailSettings{
		From:     c.FormValue("from"),
		Host:     c.FormValue("host"),
		Port:     c.FormValue("port"),
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}
	if err := h.emailServicer.SaveEmailSettings(&settings); err != nil {
		return utils.Danger(err.Error())
	}

	utils.Success(c, "Email settings saved successfully")
	settings.Password = ""
	return Render(c, pages.EmailSettingsForm(settings))
}

func (h *SettingsHandler) SendTestEmail(c echo.Context) error {
	emailSettings, err := h.emailServicer.GetEmailSettings()
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.Warning("Email settings are not properly configured")
		}
		return utils.Warning(err.Error())
	}
	if err := h.emailServicer.SendEmail([]string{emailSettings.From}, "Test email", "This is a test email"); err != nil {
		return utils.Danger("Failed to send test email, please your email settings")
	}

	return utils.Success(c, "Test email sent successfully")
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

	c.Response().Header().Add("HX-Redirect", "/auth/login")
	return nil
}

func (h *SettingsHandler) SetTheme(c echo.Context) error {
	theme := c.FormValue("theme")
	cookie := new(http.Cookie)
	cookie.Name = "theme"
	cookie.Value = theme
	cookie.Path = "/"
	c.SetCookie(cookie)

	c.Response().Header().Add("HX-Redirect", "/settings")
	return nil
}
