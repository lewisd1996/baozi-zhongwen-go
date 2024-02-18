package config

import (
	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/app"
	"github.com/lewisd1996/baozi-zhongwen/handler"
	"github.com/lewisd1996/baozi-zhongwen/middleware"
)

/* ---------------------------------- Types --------------------------------- */

type AuthRoutes struct {
	OAuthRoute            OAuthRoutes
	Login                 string
	Logout                string
	Register              string
	RegisterConfirm       string
	RegisterConfirmResend string
}

type V1ApiRoutes struct {
	Card            string
	CardEdit        string
	Cards           string
	Deck            string
	Decks           string
	Health          string
	LearningSession string
}

type APIRoutes struct {
	V1 V1ApiRoutes
}

type AppRoutes struct {
	Deck            string
	Decks           string
	Home            string
	Learn           string
	LearnSummary    string
	Login           string
	Register        string
	RegisterConfirm string
}

type OAuthRoutes struct {
	Google         string
	GoogleCallback string
}

/* --------------------------------- Routes --------------------------------- */

var APIRoute = APIRoutes{
	V1: V1ApiRoutes{
		Card:            "/decks/:deck_id/cards/:card_id",
		CardEdit:        "/decks/:deck_id/cards/:card_id/edit",
		Cards:           "/decks/:deck_id/cards",
		Deck:            "/decks/:deck_id",
		Decks:           "/decks",
		Health:          "/health",
		LearningSession: "/learn/:learning_session_id",
	},
}

var AppRoute = AppRoutes{
	Deck:            "/decks/:deck_id",
	Decks:           "/decks",
	Home:            "/",
	Learn:           "/learn",
	LearnSummary:    "/learn/:learning_session_id/summary",
	Login:           "/login",
	Register:        "/register",
	RegisterConfirm: "/register/confirm",
}

var OAuthRoute = OAuthRoutes{
	Google:         "/auth/oauth/google",
	GoogleCallback: "/auth/oauth/google/callback",
}

var AuthRoute = AuthRoutes{
	OAuthRoute:            OAuthRoute,
	Login:                 "/auth/login",
	Logout:                "/auth/logout",
	Register:              "/auth/register",
	RegisterConfirm:       "/auth/register/confirm",
	RegisterConfirmResend: "/auth/register/confirm/resend",
}

var Routes = struct {
	Auth AuthRoutes
	API  APIRoutes
	App  AppRoutes
}{
	Auth: AuthRoute,
	API:  APIRoute,
	App:  AppRoute,
}

/* ------------------------------- Add routes ------------------------------- */

func AddRoutes(e *echo.Echo, a *app.App) {
	e.Static("/assets", "assets")

	// Protected group
	protectedGroup := a.Router.Group("", func(next echo.HandlerFunc) echo.HandlerFunc {
		return middleware.AuthenticatedRouteMiddleware(next, a.Auth)
	})

	// Handlers
	CardsHandler := handler.NewCardsHandler(a)
	DecksHandler := handler.NewDecksHandler(a)
	HealthHandler := handler.NewHealthHandler(a)
	HomeHandler := handler.NewHomeHandler(a)
	LearnHandler := handler.NewLearnHandler(a)
	LoginHandler := handler.NewLoginHandler(a)
	LogoutHandler := handler.NewLogoutHandler(a)
	OAuthHandler := handler.NewOAuthHandler(a)
	RegisterHandler := handler.NewRegisterHandler(a)

	// OAuth
	a.Router.GET(Routes.Auth.OAuthRoute.Google, OAuthHandler.HandleGetGoogleLogin)
	a.Router.GET(Routes.Auth.OAuthRoute.GoogleCallback, OAuthHandler.HandleGoogleLoginCallback)
	// Login
	a.Router.POST(Routes.Auth.Login, LoginHandler.HandleLoginSubmit)
	// Logout
	a.Router.GET(Routes.Auth.Logout, LogoutHandler.HandleLogout)
	// Register
	a.Router.POST(Routes.Auth.Register, RegisterHandler.HandleRegisterSubmit)
	a.Router.POST(Routes.Auth.RegisterConfirm, RegisterHandler.HandleRegisterConfirmSubmit)
	a.Router.POST(Routes.Auth.RegisterConfirmResend, RegisterHandler.HandleRegisterConfirmResend)

	// 📡 V1 API GROUPS
	apiV1 := a.Router.Group("/v1")
	protectedGroupV1 := a.Router.Group("/v1", func(next echo.HandlerFunc) echo.HandlerFunc {
		return middleware.AuthenticatedRouteMiddleware(next, a.Auth)
	})
	// Health
	apiV1.GET(Routes.API.V1.Health, HealthHandler.HandleHealth)
	// Decks
	protectedGroupV1.POST(Routes.API.V1.Decks, DecksHandler.HandleDecksSubmit)
	protectedGroupV1.DELETE(Routes.API.V1.Deck, DecksHandler.HandleDeckDelete)
	// Cards
	protectedGroupV1.POST(Routes.API.V1.Cards, CardsHandler.HandleCardSubmit)
	// Card
	protectedGroupV1.PATCH(Routes.API.V1.Card, CardsHandler.HandlePatchCard)
	protectedGroupV1.GET(Routes.API.V1.Card, CardsHandler.HandleGetCard)
	protectedGroupV1.GET(Routes.API.V1.CardEdit, CardsHandler.HandleGetCardEdit)
	// Learning
	protectedGroupV1.POST(Routes.API.V1.LearningSession, LearnHandler.HandleLearnSessionAnswerSubmit)

	// 📱 App
	// 🔓 Unauthenticated routes
	// Login
	a.Router.GET(Routes.App.Login, LoginHandler.HandleLoginShow)
	// Register
	a.Router.GET(Routes.App.Register, RegisterHandler.HandleRegisterShow)
	a.Router.GET(Routes.App.RegisterConfirm, RegisterHandler.HandleRegisterConfirmShow)
	// 🔒 Authenticated routes
	// Home
	protectedGroup.GET(Routes.App.Home, HomeHandler.HandleHomeShow)
	// Decks
	protectedGroup.GET(Routes.App.Deck, DecksHandler.HandleDeckShow)
	protectedGroup.GET(Routes.App.Decks, DecksHandler.HandleDecksShow)
	// Learning
	protectedGroup.GET(Routes.App.Learn, LearnHandler.HandleLearnShow)
	protectedGroup.GET(Routes.App.LearnSummary, LearnHandler.HandleLearnSessionSummaryShow)
}
