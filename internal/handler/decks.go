package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/internal/app"
	"github.com/lewisd1996/baozi-zhongwen/internal/view/decks"
)

type DecksHandler struct {
	app *app.App
}

/* -------------------------------------------------------------------------- */
/*                                    Init                                    */
/* -------------------------------------------------------------------------- */

func NewDecksHandler(a *app.App) DecksHandler {
	return DecksHandler{
		app: a,
	}
}

// HandleDecksShow is a handler for GET /decks
func (h DecksHandler) HandleDecksShow(c echo.Context) error {
	userId := c.Get("user_id").(string)
	toastMessage := c.QueryParam("toast_message")
	toastType := c.QueryParam("toast_type")

	decksRes, err := h.app.Dao.GetDecksByOwnerId(userId)

	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusFound, "/404")
	}

	return Render(c, decks.Show(decksRes, userId, c.Path(), toastMessage, toastType))
}

// HandleDeckShow is a handler for GET /decks/:deck_id
func (h DecksHandler) HandleDeckShow(c echo.Context) error {
	userId := c.Get("user_id").(string)
	deckId := c.Param("deck_id")

	deckRes, err := h.app.Dao.GetDeckById(deckId)

	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusFound, "/404")
	}

	cardRes, err := h.app.Dao.GetCardsByDeckId(deckId)

	return Render(c, decks.ShowDeck(deckRes, cardRes, userId, c.Path(), "", ""))
}

// HandleDecksSubmit is a handler for POST /decks
func (h DecksHandler) HandleDecksSubmit(c echo.Context) error {
	deckName, description := c.FormValue("name"), c.FormValue("description")
	userId := c.Get("user_id").(string)

	_, err := h.app.Dao.CreateDeck(deckName, description, userId)
	if err != nil {
		log.Println(err)
		return HTML(c, decks.CreateDeckForm(err))
	}

	c.Response().Header().Set("HX-Redirect", "/decks")

	return c.NoContent(http.StatusCreated)
}

// HandleDeckDelete is a handler for DELETE /decks/:deck_id
func (h DecksHandler) HandleDeckDelete(c echo.Context) error {
	userId := c.Get("user_id").(string)
	deckId := c.Param("deck_id")
	err := h.app.Dao.DeleteDeck(deckId, userId)

	if err != nil {
		log.Println(err)
		return HTML(c, decks.CreateDeckForm(err))
	}

	c.Response().Header().Set("HX-Redirect", "/decks?toast_message=Deleted&toast_type=success")

	return c.NoContent(http.StatusCreated)
}
