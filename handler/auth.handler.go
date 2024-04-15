package handler

import (
	"crypto/rand"
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/views/auth"
	"go.uber.org/zap"
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
	fmt.Println(email)
	fmt.Println(password)
	verified, _, err := l.repo.VerifyUserWithPassword(email, password)
	if err != nil {
		return err
	}
	fmt.Println(verified)
	if verified {
		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path: "/",
			//todo: this is seven days
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}
		fmt.Println("error weee")
		sessionId, _ := generateRandomString(255)
		fmt.Println("error weee")
		err := l.repo.AddSessionId(email, sessionId)
		if err != nil {
			fmt.Println("error weee")
			return err
		}
		sess.Values["sessionId"] = sessionId
		sess.Save(c.Request(), c.Response())
		c.Response().Header().Add("HX-Redirect", "/users")
		return nil
	} else {
		//todo: show the user that the login failed -> set a header and then use alpinejs
		return Render(c, auth.Login())
	}
}

func (l LoginPageHandler) HandleRenderSingUp(c echo.Context) error {
	return Render(c, auth.Signup())
}

func (l LoginPageHandler) HandleSingUp(c echo.Context) error {
	fmt.Println("hellooo")
	email := c.FormValue("email")
	fmt.Println(email)
	password := c.FormValue("password")
	fmt.Println(password)
	username := c.FormValue("username")
	fmt.Println(username)
	ayo := c.FormValue("admin")
	fmt.Println(ayo)
	admin := c.FormValue("admin") == "admin"
	_, err := l.repo.AddUser(model.User{Email: email, Password: password, Name: username, Admin: admin})
	if err != nil {
		return err
	}
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path: "/",
		//todo: this is seven days
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
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
