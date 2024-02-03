package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/lewisd1996/baozi-zhongwen/app"
	"github.com/lewisd1996/baozi-zhongwen/handler"
	"github.com/lewisd1996/baozi-zhongwen/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
		os.Exit(1)
	}

	// Create new app
	a := app.NewApp()
	domain := os.Getenv("RAILWAY_PUBLIC_DOMAIN")

	// ⚙️ Middleware
	a.Router.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format: "time=${time_rfc3339} | method=${method} | uri=${uri} | status=${status} | host=${host}\n",
	}))
	a.Router.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(20)))
	a.Router.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "https://" + domain},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	a.Router.Use(echoMiddleware.TimeoutWithConfig(echoMiddleware.TimeoutConfig{
		Skipper:      echoMiddleware.DefaultSkipper,
		ErrorMessage: "custom timeout error message returns to client",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			fmt.Println("custom timeout error handler")
		},
		Timeout: 30 * time.Second,
	}))

	// 🗄️ Static assets
	a.Router.Static("/assets", "assets")

	// 📡 API (V1)
	v1 := a.Router.Group("/v1")
	// ├── Health
	v1.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})
	// ├── Auth
	// │   ├── Login
	LoginHandler := handler.NewLoginHandler(a)
	v1.POST("/login", LoginHandler.HandleLoginSubmit)
	// │   ├── Logout
	LogoutHandler := handler.NewLogoutHandler(a)
	v1.POST("/logout", LogoutHandler.HandleLogoutSubmit)
	// │   ├── Register
	RegisterHandler := handler.NewRegisterHandler(a)
	v1.POST("/register", RegisterHandler.HandleRegisterSubmit)
	v1.POST("/register/confirm", RegisterHandler.HandleRegisterConfirmSubmit)
	v1.POST("/register/confirm/resend", RegisterHandler.HandleRegisterConfirmResend)

	// 📱 App
	ag := a.Router.Group("", func(next echo.HandlerFunc) echo.HandlerFunc {
		return middleware.AuthenticatedRouteMiddleware(next, a.Auth)
	})
	// 🔓 Unauthenticated routes
	// ├── Auth
	// │   ├── Login
	a.Router.GET("/login", LoginHandler.HandleLoginShow)
	// │   ├── Register
	a.Router.GET("/register", RegisterHandler.HandleRegisterShow)
	a.Router.GET("/register/confirm", RegisterHandler.HandleRegisterConfirmShow)

	// 🔒 Authenticated routes
	// ├── Home
	HomeHandler := handler.NewHomeHandler()
	ag.GET("/", HomeHandler.HandleHomeShow)
	// ├── Decks
	DecksHandler := handler.NewDecksHandler()
	ag.GET("/decks", DecksHandler.HandleDecksShow)

	// Start server
	a.Router.Start(":3000")

	// Close database connection
	defer a.DB.Close()
}

// //TODO: Remove
// a.Router.GET("/users", func(c echo.Context) error {
// 	stmt := SELECT(table.User.AllColumns).FROM(table.User.Table)
// 	var res []User
// 	err := stmt.Query(a.DB, &res)
// 	if err != nil {
// 		fmt.Println(err)
// 		return c.String(500, "Error")
// 	}

// 	return c.JSON(200, res)
// })
