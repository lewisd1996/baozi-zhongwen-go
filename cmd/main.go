package main

import (
	"encoding/json"
	"os"

	"github.com/labstack/echo/v4"
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

	// Unmarshal JSON data into Vocab struct
	var entries []model.Vocab
	err = json.Unmarshal(data, &entries)
	if err != nil {
		log.Fatal(err)
	}

	// Create new app
	app := newApp(entries)

	// Register routes
	app.router.Static("/assets", "assets")
	// ├── Home
	HomeHandler := handler.NewHomeHandler(entries)
	app.router.GET("/", HomeHandler.HandleHomeShow)
	// └── Vocab
	VocabHandler := handler.NewVocabHandler(entries)
	app.router.GET("/vocab", VocabHandler.HandleVocabShow)

	// Start server
	app.router.Start(":80")

}
