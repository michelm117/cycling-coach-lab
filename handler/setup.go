package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/utils"
	"github.com/michelm117/cycling-coach-lab/views/pages"

	"go.uber.org/zap"
)

type SetupHandler struct {
	globalSettingServicer services.GlobalSettingServicer
	userServicer          services.UserServicer
	validator             utils.Validator
	logger                *zap.SugaredLogger
}

func NewSetupHandler(
	globalSettingService services.GlobalSettingServicer,
	userService services.UserServicer,
	validators utils.Validator,
	logger *zap.SugaredLogger,
) SetupHandler {
	return SetupHandler{
		globalSettingServicer: globalSettingService,
		userServicer:          userService,
		validator:             validators,
		logger:                logger,
	}
}

func (h SetupHandler) Setup(c echo.Context) error {
	if h.globalSettingServicer.IsAppInitialized() {
		return utils.Warning("App already initialized")
	}

	user, err := h.validator.CreateValidUser(
		c.FormValue("firstname"),
		c.FormValue("lastname"),
		"admin",
		c.FormValue("email"),
		c.FormValue("dateOfBirth"),
		c.FormValue("password"),
		c.FormValue("confirmPassword"),
	)
	if err != nil {
		return utils.Warning(err.Error())
	}

	_, err = h.userServicer.AddUser(*user)
	if err != nil {
		h.logger.Error(err)
		return utils.Warning("Internal server error")
	}

	err = h.globalSettingServicer.InitializeApp()
	if err != nil {
		return utils.Warning("Internal server error")
	}

	c.Response().Header().Add("HX-Redirect", "/auth/login")
	return nil
}

func (h SetupHandler) RenderSetup(c echo.Context) error {
	if !h.globalSettingServicer.IsAppInitialized() {
		return Render(c, pages.Setup(GetTheme(c)))
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/users")
}
