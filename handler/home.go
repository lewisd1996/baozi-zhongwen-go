package handler

import (
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/app"
	"github.com/lewisd1996/baozi-zhongwen/view/home"
)

type HomeHandler struct {
	app *app.App
}

/* -------------------------------------------------------------------------- */
/*                                    Init                                    */
/* -------------------------------------------------------------------------- */

func NewHomeHandler(a *app.App) HomeHandler {
	return HomeHandler{
		app: a,
	}
}

func (h HomeHandler) HandleHomeShow(c echo.Context) error {
	userId := c.Get("user_id").(string)

	var wg sync.WaitGroup
	var totalDecks, totalCards, totalCompletedLearningSessions int
	var errDeck, errCard, errSession error

	wg.Add(3)

	go func() {
		defer wg.Done()
		totalDecks, errDeck = h.app.Dao.GetUserDeckCount(userId)
	}()

	go func() {
		defer wg.Done()
		totalCards, errCard = h.app.Dao.GetUserCardCount(userId)
	}()

	go func() {
		defer wg.Done()
		totalCompletedLearningSessions, errSession = h.app.Dao.GetUserCompletedLearningSessionCount(userId)
	}()

	wg.Wait()

	if errDeck != nil || errCard != nil || errSession != nil {
		return c.Redirect(302, "/decks")
	}

	stats := []home.Stat{
		{
			Title: "Decks",
			Href:  "/decks",
			Value: totalDecks,
		},
		{
			Title: "Cards",
			Href:  "/decks",
			Value: totalCards,
		},
		{
			Title: "Completed Learning Sessions",
			Href:  "/decks",
			Value: totalCompletedLearningSessions,
		},
	}

	return Render(c, home.Show(userId, c.Path(), stats))
}
