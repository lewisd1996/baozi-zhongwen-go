package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/lewisd1996/baozi-zhongwen/handler"
	"github.com/lewisd1996/baozi-zhongwen/model"
)

type App struct {
	id      string
	entries []model.Vocab
	router  *echo.Echo
}

func newApp(entries []model.Vocab) *App {
	return &App{
		id:      "baozi-zhongwen",
		entries: entries,
		router:  echo.New(),
	}
}

func main() {
	// Load HSK JSON data into memory
	data, err := os.ReadFile("./assets/json/hsk.json")
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal HSK JSON data into Vocab struct
	var entries []model.Vocab
	err = json.Unmarshal(data, &entries)
	if err != nil {
		log.Fatal(err)
	}

	// Create new app
	app := newApp(entries)

	// Middleware
	app.router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339} | method=${method} | uri=${uri} | status=${status} | host=${host}\n",
	}))
	app.router.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	app.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://baozi-zhongwen.com/"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	app.router.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "custom timeout error message returns to client",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			fmt.Println("custom timeout error handler")
		},
		Timeout: 30 * time.Second,
	}))

	// Register routes
	app.router.Static("/assets", "assets")
	// ├── Home
	HomeHandler := handler.NewHomeHandler(entries)
	app.router.GET("/", HomeHandler.HandleHomeShow)
	// └── Vocab
	VocabHandler := handler.NewVocabHandler(entries)
	app.router.GET("/vocab", VocabHandler.HandleVocabShow)

	// Start server
	app.router.Start(":3000")

}
