package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/internal/app"
)

type LogoutHandler struct {
	app *app.App
}

/* -------------------------------------------------------------------------- */
/*                                    Init                                    */
/* -------------------------------------------------------------------------- */

func NewLogoutHandler(a *app.App) LogoutHandler {
	return LogoutHandler{app: a}
}

func (h LogoutHandler) HandleLogout(c echo.Context) error {
	h.app.Auth.Logout(c)
	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusOK)
}
