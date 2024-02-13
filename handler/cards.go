package handler

import (
	"log"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/app"
	"github.com/lewisd1996/baozi-zhongwen/view/decks"

	// dot import so that jet go code would resemble as much as native SQL
	// dot import is not mandatory
	. "github.com/go-jet/jet/v2/postgres"
	. "github.com/lewisd1996/baozi-zhongwen/sql/.jet/bz/public/model"
	"github.com/lewisd1996/baozi-zhongwen/sql/.jet/bz/public/table"
)

type CardsHandler struct {
	app *app.App
}

func NewCardsHandler(a *app.App) CardsHandler {
	return CardsHandler{
		app: a,
	}
}

// HandleCardSubmit is a handler for POST /decks/:deck_id/cards
func (h CardsHandler) HandleCardSubmit(c echo.Context) error {
	deckId := c.Param("deck_id")
	content, translation := c.FormValue("content"), c.FormValue("translation")

	// Create card in database
	card := Card{
		Content:     content,
		Translation: translation,
		DeckID:      uuid.MustParse(deckId),
	}

	stmt := table.Card.INSERT(table.Card.Content, table.Card.Translation, table.Card.DeckID).MODEL(card).RETURNING(table.Card.AllColumns)

	var cardRes Card
	err := stmt.Query(h.app.DB, &cardRes)

	if err != nil {
		log.Println(err)
		return HTML(c, decks.CreateDeckForm(err))
	}

	return HTML(c, decks.DeckCardTableRow(cardRes, "Card created", "success"))
}

// HandlePatchCard is a handler for PATCH /decks/:deck_id/cards/:card_id/
func (h CardsHandler) HandlePatchCard(c echo.Context) error {
	cardId, deckId := c.Param("card_id"), c.Param("deck_id")

	updatedCard := Card{
		Content:     c.FormValue("content"),
		Translation: c.FormValue("translation"),
	}

	println("CARD CONTENT:", updatedCard.Content)
	println("CARD TRANSLATION:", updatedCard.Translation)

	updateStmt := table.Card.UPDATE(table.Card.Content, table.Card.Translation).MODEL(updatedCard).WHERE(table.Card.ID.EQ(UUID(uuid.MustParse(cardId))).AND(table.Card.DeckID.EQ(UUID(uuid.MustParse(deckId))))).RETURNING(table.Card.AllColumns)

	var cardRes Card
	err := updateStmt.Query(h.app.DB, &cardRes)

	if err != nil {
		log.Println("ERR:", err.Error())
		return HTML(c, decks.CreateDeckForm(err))
	}

	return HTML(c, decks.DeckCardTableRow(cardRes, "Card saved", "success"))

}

// HandleGetCardEdit is a handler for GET /decks/:deck_id/cards/:card_id/
func (h CardsHandler) HandleGetCard(c echo.Context) error {
	cardId, deckId := c.Param("card_id"), c.Param("deck_id")

	stmt := table.Card.SELECT(table.Card.AllColumns).
		WHERE(
			table.Card.DeckID.EQ(UUID(uuid.MustParse(deckId))).
				AND(table.Card.ID.EQ(UUID(uuid.MustParse(cardId)))),
		).LIMIT(1)

	var cardRes Card
	err := stmt.Query(h.app.DB, &cardRes)

	if err != nil {
		println("ERR", err.Error())
		return HTML(c, decks.CreateDeckForm(err))
	}

	println("card Res:", cardRes.ID.String())

	return HTML(c, decks.DeckCardTableRow(cardRes, "", ""))
}

// HandleGetCardEdit is a handler for GET /decks/:deck_id/cards/:card_id/edit
func (h CardsHandler) HandleGetCardEdit(c echo.Context) error {
	cardId, deckId := c.Param("card_id"), c.Param("deck_id")

	stmt := table.Card.SELECT(table.Card.AllColumns).
		WHERE(
			table.Card.DeckID.EQ(UUID(uuid.MustParse(deckId))).
				AND(table.Card.ID.EQ(UUID(uuid.MustParse(cardId)))),
		).LIMIT(1)

	var cardRes Card
	err := stmt.Query(h.app.DB, &cardRes)

	if err != nil {
		println("ERR", err.Error())
		return HTML(c, decks.CreateDeckForm(err))
	}

	println("card Res:", cardRes.ID.String())

	return HTML(c, decks.EditDeckCardTableRow(cardRes))
}
