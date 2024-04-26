package handler

import (
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

func (h *SettingsHandler) RenderSettings(c echo.Context) error {
	au := c.(model.AuthenticatedContext).User
	return Render(c, pages.SettingsIndex(au))
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
