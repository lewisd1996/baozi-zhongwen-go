package config

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/lewisd1996/baozi-zhongwen/app"
)

func AddMiddleware(a *app.App, origin string) {
	a.Router.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format: "time=${time_rfc3339} | method=${method} | uri=${uri} | status=${status} | host=${host}\n",
	}))
	a.Router.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(20)))
	a.Router.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{origin},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	a.Router.Use(echoMiddleware.TimeoutWithConfig(echoMiddleware.TimeoutConfig{
		Skipper:      echoMiddleware.DefaultSkipper,
		ErrorMessage: "Timeout error",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			fmt.Println("Timeout error: ", err)
		},
		Timeout: 30 * time.Second,
	}))
}
