package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/internal/service"
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

			refreshTokenResult, err := service.RefreshToken(refreshToken.Value)
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}

			// Before dereferencing result.AccessToken and result.RefreshToken, check if they're not nil
			if refreshTokenResult != nil && refreshTokenResult.AccessToken != nil && refreshTokenResult.RefreshToken != nil {
				newToken, err := service.ValidateToken(*refreshTokenResult.AccessToken)
				if err != nil {
					return c.Redirect(http.StatusFound, "/login")
				}
				token = newToken
				service.TokenService.SetAccessTokenCookie(c, *refreshTokenResult.AccessToken)
				service.TokenService.SetRefreshTokenCookie(c, *refreshTokenResult.RefreshToken)
			} else {
				return c.Redirect(http.StatusFound, "/login")
			}
		}

		if token != nil {
			subInterface, found := token.Get("sub")
			if !found || subInterface == nil {
				return c.Redirect(http.StatusFound, "/login")
			}

			// Safely assert subInterface's type to string
			sub, ok := subInterface.(string)

			if !ok {
				return c.Redirect(http.StatusFound, "/login")
			}

			// Set the user ID in the context
			c.Set("user_id", sub)
		}

		// Token is valid, proceed with the request
		return next(c)
	}
}
