package handler

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/internal/app"
	"github.com/lewisd1996/baozi-zhongwen/internal/view/decks"
)

type CardsHandler struct {
	app *app.App
}

/* -------------------------------------------------------------------------- */
/*                                    Init                                    */
/* -------------------------------------------------------------------------- */

func NewCardsHandler(a *app.App) CardsHandler {
	return CardsHandler{
		app: a,
	}
}

// HandleCardSubmit is a handler for POST /decks/:deck_id/cards
func (h CardsHandler) HandleCardSubmit(c echo.Context) error {
	deckId := c.Param("deck_id")
	content, translation := c.FormValue("content"), c.FormValue("translation")
	cardRes, err := h.app.Dao.CreateCard(content, translation, deckId)

	if err != nil {
		log.Println(err)
		return HTML(c, decks.CreateDeckForm(err))
	}

	return HTML(c, decks.DeckCardTableRow(cardRes, "Card created", "success"))
}

// HandlePatchCard is a handler for PATCH /decks/:deck_id/cards/:card_id/
func (h CardsHandler) HandlePatchCard(c echo.Context) error {
	cardId, deckId := c.Param("card_id"), c.Param("deck_id")
	content, translation := c.FormValue("content"), c.FormValue("translation")
	cardRes, err := h.app.Dao.UpdateCard(content, translation, cardId, deckId)

	if err != nil {
		log.Println("ERR:", err.Error())
		return HTML(c, decks.CreateDeckForm(err))
	}

	return HTML(c, decks.DeckCardTableRow(cardRes, "Card saved", "success"))
}

// HandleGetCardEdit is a handler for GET /decks/:deck_id/cards/:card_id/
func (h CardsHandler) HandleGetCard(c echo.Context) error {
	cardId, deckId := c.Param("card_id"), c.Param("deck_id")
	cardRes, err := h.app.Dao.GetCardById(cardId, deckId)

	if err != nil {
		println("ERR", err.Error())
		return HTML(c, decks.CreateDeckForm(err))
	}

	return HTML(c, decks.DeckCardTableRow(cardRes, "", ""))
}

// HandleGetCardEdit is a handler for GET /decks/:deck_id/cards/:card_id/edit
func (h CardsHandler) HandleGetCardEdit(c echo.Context) error {
	cardId, deckId := c.Param("card_id"), c.Param("deck_id")
	cardRes, err := h.app.Dao.GetCardById(cardId, deckId)

	if err != nil {
		println("ERR", err.Error())
		return HTML(c, decks.CreateDeckForm(err))
	}

	return HTML(c, decks.EditDeckCardTableRow(cardRes))
}

// HandleDeleteCard is a handler for DELETE /decks/:deck_id/cards/:card_id
func (h CardsHandler) HandleDeleteCard(c echo.Context) error {
	cardId, deckId := c.Param("card_id"), c.Param("deck_id")
	err := h.app.Dao.DeleteCard(cardId, deckId)

	if err != nil {
		log.Println(err)
		return HTML(c, decks.CreateDeckForm(err))
	}

	return c.NoContent(200)
}
