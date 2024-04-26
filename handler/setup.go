package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/views/pages"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type SetupHandler struct {
	globalSettingService *services.GlobalSettingService
	userService          *services.UserService
	logger               *zap.SugaredLogger
}

func NewSetupHandler(
	globalSettingService *services.GlobalSettingService,
	userService *services.UserService,
	logger *zap.SugaredLogger,
) SetupHandler {
	return SetupHandler{
		globalSettingService: globalSettingService,
		userService:          userService,
		logger:               logger,
	}
}

func (h SetupHandler) Setup(c echo.Context) error {
	if h.globalSettingService.IsAppInitialized() {
		return Warning("App already initialized")
	}

	firstname := c.FormValue("firstname")
	lastname := c.FormValue("lastname")
	email := c.FormValue("email")

	dateOfBirthStr := c.FormValue("dateOfBirth")
	dateOfBirth, err := time.Parse("2006-01-02", dateOfBirthStr)
	if err != nil {
		return Warning("Invalid date of birth")
	}

	password := c.FormValue("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Warning("Internal server error")
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
	_, err = h.userService.AddUser(u)
	if err != nil {
		return Warning("Internal server error")
	}

	h.globalSettingService.InitializeApp()

	c.Response().Header().Add("HX-Redirect", "/auth/login")
	return nil
}

func (h SetupHandler) RenderSetup(c echo.Context) error {
	if !h.globalSettingService.IsAppInitialized() {
		return Render(c, pages.Setup())
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/users")
}
