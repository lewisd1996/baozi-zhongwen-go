package service

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type TokenService struct{}

/* -------------------------------------------------------------------------- */
/*                                    Init                                    */
/* -------------------------------------------------------------------------- */

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (service *TokenService) SetAccessTokenCookie(c echo.Context, accessToken string) {
	accessTokenCookie := http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",                 // available to all paths
		Domain:   os.Getenv("DOMAIN"), // set your domain here (e.g. localhost, myapp.com, etc.)
		HttpOnly: true,
		Secure:   true, // consider the environment, as mentioned above
		// SameSite: http.SameSiteLaxMode, // Uncomment if necessary
	}

	c.SetCookie(&accessTokenCookie)
}

func (service *TokenService) SetRefreshTokenCookie(c echo.Context, refreshToken string) {
	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken, // Set your refresh token here
		Path:     "/",
		Domain:   os.Getenv("DOMAIN"),                  // set your domain here (e.g. localhost, myapp.com, etc.)
		Expires:  time.Now().Add(365 * 24 * time.Hour), // 1 year, adjust based on your requirements
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(&refreshTokenCookie)
}

func (service *TokenService) DeleteAccessTokenCookie(c echo.Context) {
	deletedAccessTokenCookie := http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	c.SetCookie(&deletedAccessTokenCookie)
}

func (service *TokenService) DeleteRefreshTokenCookie(c echo.Context) {
	deletedRefreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	c.SetCookie(&deletedRefreshTokenCookie)
}
