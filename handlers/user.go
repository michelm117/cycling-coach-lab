package handlers

import (
	"github.com/labstack/echo/v4"
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
	users := []models.User{
		{
			Email: "mar@mas.lu",
		},
		{
			Email: "sadasd@adsasd.de",
		},
	}
	return render(c, user.ShowUsers(users))
}
