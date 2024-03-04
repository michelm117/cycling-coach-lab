package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/michelm117/cycling-coach-lab/models"
	"github.com/michelm117/cycling-coach-lab/repositories"
	"github.com/michelm117/cycling-coach-lab/views/components"
	"github.com/michelm117/cycling-coach-lab/views/user"
)

type UserHandler struct {
	repo *repositories.UserRepository
}

func NewUserHandler(repo *repositories.UserRepository) UserHandler {
	return UserHandler{repo: repo}
}

func (h UserHandler) HandleDeleteUser(c echo.Context) error {
	email := c.ParamValues()
	emailOfUser := strings.Replace(email[0], "%40", "@", -1)
	userToBeDeleted, err := h.repo.GetByEmail(emailOfUser)
	if err != nil {
		return err
	}
	h.repo.DeleteUser(*userToBeDeleted)
	users, _ := h.repo.GetAllUsers()
	for _, t := range users {
		println(t.Email)
		println(t.Name)
	}

	return render(c, user.ShowUsers(users))
}

func (h UserHandler) HandlerShowUserById(c echo.Context) error {
	u := models.User{
		Email: "a@gg.com",
	}
	return render(c, user.ShowUser(u))
}

func (h UserHandler) HandlerGetUser(c echo.Context) error {
	keyword := c.FormValue("search")

	users, _ := h.repo.SearchForUser(keyword)
	return render(c, user.ShowUsers(users))
}

func (h UserHandler) HandlerShowUsers(c echo.Context) error {
	users, err := h.repo.GetAllUsers()
	if err != nil {
		fmt.Println("error when looking for all users:" + err.Error())
	}
	return render(c, user.ShowUsers(users))
}

func (h UserHandler) HandlerAddUser(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")

	newUser := models.User{
		Name:  name,
		Email: email,
	}
	_, err := h.repo.AddUser(newUser)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return render(c, components.EmailTaken(newUser))
		}
		fmt.Println("error when adding user:" + err.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	users, _ := h.repo.GetAllUsers()
	return render(c, user.ShowUsers(users))
}
