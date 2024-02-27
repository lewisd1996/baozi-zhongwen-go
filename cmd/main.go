package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"

	"github.com/lewisd1996/baozi-zhongwen/internal/app"
	"github.com/lewisd1996/baozi-zhongwen/internal/config"
)

// type Config struct {
// 	ctx    context.Context
// 	getenv func(string) string
// 	stdout io.Writer
// 	stderr io.Writer
// 	args   []string
// }

func run() error {
	origin := os.Getenv("URL")

	// 🚀 Initialize app
	a := app.NewApp()

	// ⚙️ Middleware
	config.AddMiddleware(a, origin)

	// 📡 Routes
	config.AddRoutes(a.Router, a)

	go func() {
		if err := a.Router.Start(":3000"); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "server failed to start: %s\n", err)
			os.Exit(1)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.Router.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "server shutdown failed: %s\n", err)
	}

	defer a.Dao.DB.Close()
	return nil
}

func main() {
	appEnv := os.Getenv("APP_ENV")

	if appEnv == "development" {
		if err := godotenv.Load(); err != nil {
			fmt.Println("No .env file found")
			os.Exit(1)
		}
	}

	// config := Config{
	// 	ctx:    context.Background(),
	// 	getenv: os.Getenv,
	// 	stdout: os.Stdout,
	// 	stderr: os.Stderr,
	// 	args:   os.Args,
	// }

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
