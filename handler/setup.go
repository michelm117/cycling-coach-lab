package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/utils"
	"github.com/michelm117/cycling-coach-lab/views/pages"

	"go.uber.org/zap"
)

type SetupHandler struct {
	globalSettingServicer services.GlobalSettingServicer
	userServicer          services.UserServicer
	cryptoer              utils.Cryptoer
	logger                *zap.SugaredLogger
}

func NewSetupHandler(
	globalSettingService services.GlobalSettingServicer,
	userService services.UserServicer,
	cryptoer utils.Cryptoer,
	logger *zap.SugaredLogger,
) SetupHandler {
	return SetupHandler{
		globalSettingServicer: globalSettingService,
		userServicer:          userService,
		cryptoer:              cryptoer,
		logger:                logger,
	}
}

func (h SetupHandler) Setup(c echo.Context) error {
	if h.globalSettingServicer.IsAppInitialized() {
		return utils.Warning("App already initialized")
	}

	firstname := c.FormValue("firstname")
	lastname := c.FormValue("lastname")
	email := c.FormValue("email")

	dateOfBirthStr := c.FormValue("dateOfBirth")
	dateOfBirth, err := time.Parse("2006-01-02", dateOfBirthStr)
	if err != nil {
		return utils.Warning("Invalid date of birth")
	}

	password := c.FormValue("password")
	hashedPassword, err := h.cryptoer.GenerateFromPassword([]byte(password))
	if err != nil {
		return utils.Warning("Internal server error")
	}

	u := model.User{
		Firstname:    firstname,
		Lastname:     lastname,
		Email:        email,
		DateOfBirth:  dateOfBirth,
		Role:         "admin",
		Status:       "active",
		PasswordHash: string(hashedPassword),
	}

	_, err = h.userServicer.AddUser(u)
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
		return Render(c, pages.Setup())
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/users")
}
