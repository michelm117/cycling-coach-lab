package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/michelm117/cycling-coach-lab/db"
	"github.com/michelm117/cycling-coach-lab/models"
	"github.com/michelm117/cycling-coach-lab/views/user"
)

type UserHandler struct{}

func (h UserHandler) HandlerShowUserById(c echo.Context) error {
	u := models.User{
		Email: "a@gg.com",
	}
	return render(c, user.ShowUser(u))
}

func (h UserHandler) HandlerShowUsers(c echo.Context) error {
	fmt.Println("handle all users")
	users, err := db.GetAllUsers()
	if err != nil {
		fmt.Println("error when looking for all users:" + err.Error())
	}

	if len(users) == 0 {
		println("no users in db")
		users = append(users, models.User{
			Name:  "name",
			Email: "email",
		})
	}
	println("we got till here")
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
