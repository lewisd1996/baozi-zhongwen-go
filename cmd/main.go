package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"

	"github.com/lewisd1996/baozi-zhongwen/app"
	"github.com/lewisd1996/baozi-zhongwen/config"
)

func run(
	ctx context.Context,
	getenv func(string) string,
	stdout, stderr io.Writer,
	args []string,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	domain := os.Getenv("DOMAIN")

	// 🚀 Initialize app
	a := app.NewApp()

	// ⚙️ Middleware
	config.AddMiddleware(a, domain)

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

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.Router.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "server shutdown failed: %s\n", err)
	}

	defer a.DB.Close()
	return nil
}

func main() {
	ctx := context.Background()

	appEnv := os.Getenv("APP_ENV")

	if appEnv == "development" {
		if err := godotenv.Load(); err != nil {
			fmt.Println("No .env file found")
			os.Exit(1)
		}
	}

	if err := run(
		ctx,
		os.Getenv,
		os.Stdout,
		os.Stderr,
		os.Args,
	); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
