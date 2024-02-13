package config

import (
	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/app"
	"github.com/lewisd1996/baozi-zhongwen/handler"
	"github.com/lewisd1996/baozi-zhongwen/middleware"
)

func AddRoutes(e *echo.Echo, a *app.App) {
	e.Static("/assets", "assets")

	ag := a.Router.Group("", func(next echo.HandlerFunc) echo.HandlerFunc {
		return middleware.AuthenticatedRouteMiddleware(next, a.Auth)
	})

	// GROUPS
	v1 := a.Router.Group("/v1")
	v1ag := ag.Group("/v1")

	// Handlers
	CardsHandler := handler.NewCardsHandler(a)
	DecksHandler := handler.NewDecksHandler(a)
	HomeHandler := handler.NewHomeHandler()
	LearnHandler := handler.NewLearnHandler(a)
	LoginHandler := handler.NewLoginHandler(a)
	LogoutHandler := handler.NewLogoutHandler(a)
	RegisterHandler := handler.NewRegisterHandler(a)

	// 📡 API GROUPS (V1)
	// ├── Health
	v1.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})
	// ├── Auth
	// │   ├── Login
	v1.POST("/login", LoginHandler.HandleLoginSubmit)
	// │   ├── Logout
	v1.POST("/logout", LogoutHandler.HandleLogoutSubmit)
	// │   ├── Register
	v1.POST("/register", RegisterHandler.HandleRegisterSubmit)
	v1.POST("/register/confirm", RegisterHandler.HandleRegisterConfirmSubmit)
	v1.POST("/register/confirm/resend", RegisterHandler.HandleRegisterConfirmResend)
	// ├── Decks
	v1ag.POST("/decks", DecksHandler.HandleDecksSubmit)
	v1ag.DELETE("/decks/:deck_id", DecksHandler.HandleDeckDelete)
	// ├── Cards
	v1ag.POST("/decks/:deck_id/cards", CardsHandler.HandleCardSubmit)
	v1ag.PATCH("/decks/:deck_id/cards/:card_id", CardsHandler.HandlePatchCard)
	v1ag.GET("/decks/:deck_id/cards/:card_id", CardsHandler.HandleGetCard)
	v1ag.GET("/decks/:deck_id/cards/:card_id/edit", CardsHandler.HandleGetCardEdit)
	// ├── Learning
	v1ag.POST("/learn/:learning_session_id", LearnHandler.HandleLearnSessionAnswerSubmit)

	// 📱 App
	// 🔓 Unauthenticated routes
	// ├── Auth
	// │   ├── Login
	a.Router.GET("/login", LoginHandler.HandleLoginShow)
	// │   ├── Register
	a.Router.GET("/register", RegisterHandler.HandleRegisterShow)
	a.Router.GET("/register/confirm", RegisterHandler.HandleRegisterConfirmShow)
	// 🔒 Authenticated routes
	// ├── Home
	ag.GET("/", HomeHandler.HandleHomeShow)
	// ├── Decks
	ag.GET("/decks", DecksHandler.HandleDecksShow)
	ag.GET("/decks/:deck_id", DecksHandler.HandleDeckShow)
	// ├── Learning
	ag.GET("/learn", LearnHandler.HandleLearnShow)
	ag.GET("/learn/:learning_session_id/summary", LearnHandler.HandleLearnSessionSummaryShow)
}
