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
			return c.Redirect(http.StatusFound, "/login")
		}
		tokenString := cookie.Value

		// Use your JWK set to validate the token
		token, err := service.ValidateToken(tokenString)

		if err != nil {
			// If access token is expired, try to refresh it
			refreshToken, err := c.Cookie("refresh_token")
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}

			result, err := service.RefreshToken(refreshToken.Value)
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}
			setTokenCookies(c, result)
			newToken, err := service.ValidateToken(*result.AccessToken)
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}
			token = newToken
		}

		if token != nil {
			subInterface, found := token.Get("sub")
			if !found || subInterface == nil {
				// Handle missing or nil sub claim
				return c.Redirect(http.StatusFound, "/login")
			}

			// Safely assert subInterface's type to string
			sub, ok := subInterface.(string)

			if !ok {
				// Handle case where sub is not a string
				return c.Redirect(http.StatusFound, "/login")
			}

			// Set the user ID in the context
			c.Set("user_id", sub)
		}

		// Token is valid, proceed with the request
		return next(c)
	}
}

/* ---------------------------- Utility functions --------------------------- */

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
