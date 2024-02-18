package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/michelm117/cycling-coach-lab/db"
	"github.com/michelm117/cycling-coach-lab/models"
	"github.com/michelm117/cycling-coach-lab/repositories"
	"github.com/michelm117/cycling-coach-lab/views/user"
)

type UserHandler struct {
	repo *repositories.UserRepository
}

func (h UserHandler) HandlerShowUserById(c echo.Context) error {
	u := models.User{
		Email: "a@gg.com",
	}
	return render(c, user.ShowUser(u))
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
	println(newUser.Name)
	println(newUser.Email)
	err := db.AddUser(newUser)
	if err != nil {
		fmt.Println("error when adding user:" + err.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	users, err := db.GetAllUsers()
	return render(c, user.ShowUsers(users))
}
