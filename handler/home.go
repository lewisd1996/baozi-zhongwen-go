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
	var err error

	wg.Add(3)

	go func() {
		defer wg.Done()
		totalDecks, err = h.app.Dao.GetUserDeckCount(userId)
		if err != nil {
			return
		}
	}()

	go func() {
		defer wg.Done()
		totalCards, err = h.app.Dao.GetUserCardCount(userId)
		if err != nil {
			return
		}
	}()

	go func() {
		defer wg.Done()
		totalCompletedLearningSessions, err = h.app.Dao.GetUserCompletedLearningSessionCount(userId)
		if err != nil {
			return
		}
	}()

	wg.Wait()

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
