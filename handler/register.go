package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/app"
	"github.com/lewisd1996/baozi-zhongwen/view/auth/confirm"
	"github.com/lewisd1996/baozi-zhongwen/view/auth/register"
)

type RegisterHandler struct {
	app *app.App
}

func NewRegisterHandler(a *app.App) RegisterHandler {
	return RegisterHandler{app: a}
}

func (h RegisterHandler) HandleRegisterShow(c echo.Context) error {
	return Render(c, register.Show(c.Path()))
}

func (h RegisterHandler) HandleRegisterConfirmShow(c echo.Context) error {
	username := c.QueryParam("username")
	return Render(c, confirm.Show(username, c.Path(), false))
}

func (h RegisterHandler) HandleRegisterSubmit(c echo.Context) error {
	username, password := c.FormValue("username"), c.FormValue("password")
	authResult, err := h.app.Auth.Register(username, password)

	if err != nil {
		log.Println(err)
		return HTML(c, register.RegisterForm(err))
	}

	// Create user in database
	userSub := *authResult.UserSub
	userId, err := uuid.Parse(userSub)
	err = h.app.Dao.CreateUser(username, userId)

	if err != nil {
		log.Println(err)
		return HTML(c, register.RegisterForm(fmt.Errorf("Failed to create user")))
	}

	if authResult.UserConfirmed == nil || !*authResult.UserConfirmed {
		c.Response().Header().Set("HX-Redirect", "/register/confirm?username="+username)
		return c.NoContent(http.StatusOK)
	}

	return c.NoContent(http.StatusOK)
}

func (h RegisterHandler) HandleRegisterConfirmSubmit(c echo.Context) error {
	username := c.FormValue("username")
	code := c.FormValue("code")

	err := h.app.Auth.Confirm(username, code)

	if err != nil {
		log.Println(err)
		return HTML(c, confirm.ConfirmForm(username, err, false))
	}

	c.Response().Header().Set("HX-Redirect", "/login")

	return c.NoContent(http.StatusOK)
}

func (h RegisterHandler) HandleRegisterConfirmResend(c echo.Context) error {
	username := c.FormValue("username")
	err := h.app.Auth.ResendConfirmationCode(username)
	if err != nil {
		return HTML(c, confirm.ConfirmForm(username, err, true))
	}

	return c.NoContent(http.StatusOK)
}
