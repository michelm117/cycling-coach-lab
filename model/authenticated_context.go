package model

import (
	"github.com/labstack/echo/v4"
)

type AuthenticatedContext struct {
	echo.Context
	User *User
}
