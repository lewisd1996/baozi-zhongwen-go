package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/app"
	"github.com/lewisd1996/baozi-zhongwen/view/auth/login"
)

type LoginHandler struct {
	app *app.App
}

func NewLoginHandler(a *app.App) LoginHandler {
	return LoginHandler{app: a}
}

func (h LoginHandler) HandleLoginShow(c echo.Context) error {
	err := c.QueryParam("error")
	var errorMessage string
	if err != "" {
		if err == "not_registered" {
			errorMessage = "You are not registered. Please register first."
		} else {
			errorMessage = fmt.Sprintf("Error: %s", err)
		}
	}
	return Render(c, login.Show(errorMessage, c.Path()))
}

func (h LoginHandler) HandleLoginSubmit(c echo.Context) error {
	username, password := c.FormValue("username"), c.FormValue("password")

	err := h.app.Auth.LoginWithUsernamePassword(c, username, password)

	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Redirect", "/")

	return c.NoContent(http.StatusOK)
}
