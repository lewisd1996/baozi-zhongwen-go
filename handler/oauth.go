package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/app"
)

type OAuthHandler struct {
	app *app.App
}

/* -------------------------------------------------------------------------- */
/*                                    Init                                    */
/* -------------------------------------------------------------------------- */

func NewOAuthHandler(a *app.App) OAuthHandler {
	return OAuthHandler{app: a}
}

/* --------------------------------- Google --------------------------------- */

func (h OAuthHandler) HandleGetGoogleLogin(c echo.Context) error {
	c.Response().Header().Set("HX-Redirect", h.app.Auth.GetGoogleLoginURL())
	println(h.app.Auth.GetGoogleLoginURL())
	return c.NoContent(302)
}

func (h OAuthHandler) HandleGoogleLoginCallback(c echo.Context) error {
	code := c.QueryParam("code")

	if code == "" {
		return c.Redirect(302, "/auth/login?error=google_oauth_error")
	}

	err := h.app.Auth.SignInWithGoogle(c, code)

	if err != nil {
		return c.Redirect(302, "/auth/login?error=google_oauth_error")
	}

	return c.Redirect(302, "/")
}
