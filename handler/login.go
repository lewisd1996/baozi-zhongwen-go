package handler

import (
	"log"
	"net/http"
	"net/url"
	"time"

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
	print("LOGIN PATH:", c.Path())
	return Render(c, login.Show(c.Path()))
}

func (h LoginHandler) HandleLoginSubmit(c echo.Context) error {
	username, password := c.FormValue("username"), c.FormValue("password")
	authResult, err := h.app.Auth.Login(username, password)

	if err != nil {
		log.Println("[LOGIN ERROR]:", err.Error())

		if err.Error() == "UserNotConfirmedException: User is not confirmed." {
			encodedUsername := url.QueryEscape(username)
			c.Response().Header().Set("HX-Redirect", "/register/confirm?username="+encodedUsername)
			return c.NoContent(http.StatusOK)
		}

		return HTML(c, login.LoginForm(err))
	}

	accessTokenCookie := http.Cookie{
		Name:     "access_token",
		Value:    *authResult.AccessToken,
		Path:     "/",         // available to all paths
		Domain:   "localhost", // adjust as per your domain, might omit for localhost
		HttpOnly: true,
		Secure:   true, // consider the environment, as mentioned above
		// SameSite: http.SameSiteLaxMode, // Uncomment if necessary
	}

	c.SetCookie(&accessTokenCookie)

	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    *authResult.RefreshToken, // Set your refresh token here
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(365 * 24 * time.Hour), // 1 year, adjust based on your requirements
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(&refreshTokenCookie)

	c.Response().Header().Set("HX-Redirect", "/")

	return c.NoContent(http.StatusOK)
}
