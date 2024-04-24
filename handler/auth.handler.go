package handler

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/views/auth"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type LoginPageHandler struct {
	repo   *services.UserService
	logger *zap.SugaredLogger
}

func NewLoginPageHandler(
	repo *services.UserService,
	logger *zap.SugaredLogger,
) LoginPageHandler {
	return LoginPageHandler{repo: repo}
}

func (l LoginPageHandler) HandleRenderLogin(c echo.Context) error {
	return Render(c, auth.Login())
}

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLength := len(charset)

	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	for i, b := range randomBytes {
		randomBytes[i] = charset[int(b)%charsetLength]
	}

	return string(randomBytes), nil
}

func (l LoginPageHandler) HandleLogin(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	user, err := l.repo.GetByEmail(email)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return Warning("Invalid credentials")
	}
	fmt.Println(string(hashedPassword))
	fmt.Println(user.PasswordHash)
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path: "/",
		// todo: this is seven days
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	// Todo: use uuid by postgres
	sessionId, _ := generateRandomString(255)
	err = l.repo.AddSessionId(email, sessionId)
	if err != nil {
		return err
	}
	sess.Values["sessionId"] = sessionId
	sess.Save(c.Request(), c.Response())
	c.Response().Header().Add("HX-Redirect", "/users")
	return nil
}

func (l LoginPageHandler) HandleRenderSingUp(c echo.Context) error {
	return Render(c, auth.Signup())
}

func (l LoginPageHandler) HandleSingUp(c echo.Context) error {
	firstname := c.FormValue("firstname")
	lastname := c.FormValue("lastname")
	email := c.FormValue("email")
	role := c.FormValue("role")

	dateOfBirthStr := c.FormValue("dateOfBirth")
	dateOfBirth, err := time.Parse("2006-01-02", dateOfBirthStr)

	password := c.FormValue("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Warning("Internal Server Error")
	}

	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path: "/",
		// todo: this is seven days
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sessionId, _ := generateRandomString(255)
	newUser := model.User{
		Firstname:    firstname,
		Lastname:     lastname,
		DateOfBirth:  dateOfBirth,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		Status:       "active",
		SessionId:    sessionId,
	}
	_, err = l.repo.AddUser(newUser)
	if err != nil {
		return err
	}
	sess.Values["sessionId"] = sessionId
	sess.Save(c.Request(), c.Response())
	fmt.Println("redirect!")
	c.Response().Header().Add("HX-Redirect", "/users")
	return nil
}
