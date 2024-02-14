package handler

import (
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

	usersDeckCount, err := h.app.Dao.GetUserDeckCount(userId)
	usersCardCount, err := h.app.Dao.GetUserCardCount(userId)
	usersCompletedLearningSessionCount, err := h.app.Dao.GetUserCompletedLearningSessionCount(userId)

	if err != nil {
		return err
	}

	var stats []home.Stat

	stats = append(stats, home.Stat{
		Title: "Decks",
		Href:  "/decks",
		Value: usersDeckCount,
	})
	stats = append(stats, home.Stat{
		Title: "Cards",
		Href:  "/decks",
		Value: usersCardCount,
	})
	stats = append(stats, home.Stat{
		Title: "Completed Learning Sessions",
		Href:  "/decks",
		Value: usersCompletedLearningSessionCount,
	})

	return Render(c, home.Show(userId, c.Path(), stats))
}
