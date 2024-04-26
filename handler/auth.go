package handler

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/views/pages"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionMaxAgeSeconds = 86400 * 7
	sessionCookieName    = "cycling-coach-lab"
)

type AuthHandler struct {
	userService          *services.UserService
	sessionService       *services.SessionService
	globalSettingService *services.GlobalSettingService
	logger               *zap.SugaredLogger
}

func NewAuthHandler(
	userService *services.UserService,
	sessionService *services.SessionService,
	globalSettingService *services.GlobalSettingService,
	logger *zap.SugaredLogger,
) AuthHandler {
	return AuthHandler{
		userService:          userService,
		sessionService:       sessionService,
		globalSettingService: globalSettingService,
		logger:               logger,
	}
}

func (h AuthHandler) RenderLogin(c echo.Context) error {
	if !h.globalSettingService.IsAppInitialized() {
		return c.Redirect(http.StatusTemporaryRedirect, "/setup")
	}

	return Render(c, pages.Login())
}

func (h AuthHandler) Login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Retrieve user by email
	user, err := h.userService.GetByEmail(email)
	if err != nil {
		return Warning("Invalid credentials")
	}

	// Compare password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return Warning("Invalid credentials")
	}

	// Get the session
	browserSession, _ := session.Get(sessionCookieName, c)

	// Configure session options
	browserSession.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionMaxAgeSeconds,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") != "development",
	}

	// Save session ID
	sessionID, err := h.sessionService.SaveSession(user.ID)
	if err != nil {
		return Warning("Invalid credentials")
	}
	browserSession.Values["sessionId"] = sessionID

	if err := browserSession.Save(c.Request(), c.Response()); err != nil {
		return Warning("Invalid credentials")
	}

	c.Response().Header().Add("HX-Redirect", "/users")

	return nil
}
