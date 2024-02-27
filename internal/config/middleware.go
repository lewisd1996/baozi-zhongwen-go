package config

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/lewisd1996/baozi-zhongwen/internal/app"
	"github.com/rs/zerolog"
)

func AddMiddleware(a *app.App, origin string) {
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

	buildInfo, _ := debug.ReadBuildInfo()
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()

	a.Router.Use(echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")
			return nil
		},
	}))
}
