package handler

import (
	"fmt"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/services"
	"github.com/michelm117/cycling-coach-lab/utils"
	"github.com/michelm117/cycling-coach-lab/views/pages"
)

type UserManagementHandler struct {
	userServicer  services.UserServicer
	emailServicer services.EmailServicer
	validator     utils.Validator
	logger        *zap.SugaredLogger
}

func NewUserManagementHandler(
	userServicer services.UserServicer,
	emailServicer services.EmailServicer,
	validator utils.Validator,
	logger *zap.SugaredLogger,
) UserManagementHandler {
	return UserManagementHandler{userServicer: userServicer, emailServicer: emailServicer, validator: validator, logger: logger}
}

func (h UserManagementHandler) RenderUserManagementPage(c echo.Context) error {
	au := c.(model.AuthenticatedContext).User

	users, err := h.userServicer.GetAllUsers()
	if err != nil {
		return utils.Warning("Could not retrieve users")
	}
	return Render(c, pages.UserManagementPage(au, GetTheme(c), users))
}

func (h UserManagementHandler) RenderUserManagementView(c echo.Context) error {
	au := c.(model.AuthenticatedContext).User

	users, err := h.userServicer.GetAllUsers()
	if err != nil {
		return utils.Warning("Could not retrieve users")
	}
	return Render(c, pages.UserManagementView(au, users))
}

func (h UserManagementHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.Warning("Invalid user id")
	}

	if err := h.userServicer.DeleteUser(id); err != nil {
		return utils.Warning(fmt.Sprintf("Could not delete user with id '%v'", id))
	}

	return utils.Success(c, "User deleted successfully")
}

func (h UserManagementHandler) AddUser(c echo.Context) error {
	user, err := h.validator.CreateValidUser(
		c.FormValue("firstname"),
		c.FormValue("lastname"),
		c.FormValue("role"),
		c.FormValue("email"),
		c.FormValue("dateOfBirth"),
		c.FormValue("password"),
		c.FormValue("confirmPassword"),
	)
	if err != nil {
		return utils.Warning(err.Error())
	}

	if _, err := h.userServicer.AddUser(*user); err != nil {
		return utils.Warning("Could not add user")
	}

	utils.Success(c, fmt.Sprintf("User '%s' added successfully", user.Email))
	return Render(c, pages.AddUserResponse(user))
}

func (h UserManagementHandler) InviteUser(c echo.Context) error {
	email := c.FormValue("email")
	if err := h.validator.ValidateEmail(email); err != nil {
		return utils.Warning("Invalid email")
	}

	role := c.FormValue("role")
	if err := h.validator.ValidateRole(role); err != nil {
		return utils.Warning("Invalid role")
	}

	token := ""
	err := h.emailServicer.SendEmailWithTemplate(
		[]string{email},
		"Registrationseinladung - Cycling Coach Lab",
		"email/registration-invite.tpl",
		struct {
			Domain string
			Token  string
		}{
			Domain: os.Getenv("DOMAIN"),
			Token:  token,
		},
	)
	if err != nil {
		return utils.Warning("Could not send invitation email")
	}
	return utils.Success(c, fmt.Sprintf("User '%s' invited successfully", email))
}
