package middleware

import (
	"net/http"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/service"
)

func AuthenticatedRouteMiddleware(next echo.HandlerFunc, service *service.AuthService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract the token from the cookie
		cookie, err := c.Cookie("access_token") // Use the correct cookie name

		if err != nil {
			// Handle missing cookie (token)
			return echo.NewHTTPError(http.StatusUnauthorized, "missing token cookie")
		}
		tokenString := cookie.Value

		// Use your JWK set to validate the token
		_, err = service.ValidateToken(tokenString)
		if err != nil {
			// If access token is expired, try to refresh it
			refreshToken, err := c.Cookie("refresh_token")
			if err != nil {
				// Handle missing refresh token
				return echo.NewHTTPError(http.StatusUnauthorized, "missing refresh token")
			}
			result, err := service.RefreshToken(refreshToken.Value)
			if err != nil {
				// Handle error in refreshing tokens
				return echo.NewHTTPError(http.StatusUnauthorized, "unable to refresh tokens")
			}

			setTokenCookies(c, result)
		}

		// Token is valid, proceed with the request
		return next(c)
	}
}

// Utility functions

func setTokenCookies(c echo.Context, tokens *cognitoidentityprovider.AuthenticationResultType) {
	// Set access token cookie
	accessTokenCookie := http.Cookie{
		Name:     "access_token",
		Value:    *tokens.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	c.SetCookie(&accessTokenCookie)

	// Set refresh token cookie, if you receive a new one
	if tokens.RefreshToken != nil {
		refreshTokenCookie := http.Cookie{
			Name:     "refresh_token",
			Value:    *tokens.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		}
		c.SetCookie(&refreshTokenCookie)
	}
}
