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
	url := h.app.Auth.GetGoogleLoginURL()
	if url == "" {
		return c.JSON(500, "error getting google login url")
	}
	c.Response().Header().Set("HX-Redirect", url)
	return c.NoContent(302)
}

func (h OAuthHandler) HandleGetGoogleRegister(c echo.Context) error {
	url := h.app.Auth.GetGoogleRegisterURL()
	if url == "" {
		return c.JSON(500, "error getting google register url")
	}
	c.Response().Header().Set("HX-Redirect", url)
	return c.NoContent(302)
}

func (h OAuthHandler) HandleGoogleLoginCallback(c echo.Context) error {
	println("USER LOGGING IN WITH GOOGLE")
	code := c.QueryParam("code")

	if code == "" {
		return c.Redirect(302, "/auth/login?error=google_oauth_error")
	}

	err := h.app.Auth.SignInWithGoogle(c, code, h.app.Dao)

	if err != nil {
		if err.Error() == "user is not registered" {
			return c.Redirect(302, "/login?error=not_registered")
		}
	}

	return c.Redirect(302, "/")
}

func (h OAuthHandler) HandleGoogleRegisterCallback(c echo.Context) error {
	println("[HandleGoogleRegisterCallback]: USER REGISTERING WITH GOOGLE")
	code := c.QueryParam("code")

	if code == "" {
		println("[HandleGoogleRegisterCallback]: code is empty")
		return c.Redirect(302, "/auth/register?error=google_oauth_error")
	}

	err := h.app.Auth.RegisterWithGoogle(c, code, h.app.Dao)

	if err != nil {
		println("[HandleGoogleRegisterCallback]: error:", err.Error())
		return c.Redirect(302, "/auth/register?error=google_oauth_error")
	}

	return c.Redirect(302, "/")
}
