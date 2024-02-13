package handler

import (
	"log"
	"net/http"

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

type DecksHandler struct {
	app *app.App
}

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

	stmt := table.Deck.SELECT(table.Deck.AllColumns).WHERE(table.Deck.OwnerID.EQ(UUID(uuid.MustParse(userId)))).ORDER_BY(table.Deck.UpdatedAt)
	var res []Deck
	err := stmt.Query(h.app.DB, &res)

	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusFound, "/404")
	}

	return Render(c, decks.Show(res, userId, c.Path(), toastMessage, toastType))
}

// HandleDeckShow is a handler for GET /decks/:deck_id
func (h DecksHandler) HandleDeckShow(c echo.Context) error {
	userId := c.Get("user_id").(string)
	deckId := c.Param("deck_id")

	deckStmt := table.Deck.SELECT(table.Deck.AllColumns).WHERE(table.Deck.ID.EQ(UUID(uuid.MustParse(deckId)))).LIMIT(1)
	var deckRes Deck
	err := deckStmt.Query(h.app.DB, &deckRes)

	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusFound, "/404")
	}

	cardStmt := table.Card.SELECT(table.Card.AllColumns).WHERE(table.Card.DeckID.EQ(UUID(uuid.MustParse(deckId)))).ORDER_BY(table.Card.UpdatedAt.DESC())
	var cardRes []Card
	err = cardStmt.Query(h.app.DB, &cardRes)

	return Render(c, decks.ShowDeck(deckRes, cardRes, userId, c.Path(), "", ""))
}

// HandleDecksSubmit is a handler for POST /decks
func (h DecksHandler) HandleDecksSubmit(c echo.Context) error {
	deckName, description := c.FormValue("name"), c.FormValue("description")

	userId := c.Get("user_id").(string)

	parsedUserId, err := uuid.Parse(userId)

	if err != nil {
		log.Println(err)
		return HTML(c, decks.CreateDeckForm(err))
	}

	// Create deck in database
	deck := Deck{
		Name:        deckName,
		Description: description,
		OwnerID:     parsedUserId,
	}

	stmt := table.Deck.INSERT(table.Deck.OwnerID, table.Deck.Name, table.Deck.Description).MODEL(deck).RETURNING(table.Deck.AllColumns)
	_, err = stmt.Exec(h.app.DB)

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

	stmt := table.Deck.DELETE().
		WHERE(table.Deck.OwnerID.EQ(UUID(uuid.MustParse(userId))).AND(table.Deck.ID.EQ(UUID(uuid.MustParse(deckId))))).
		RETURNING(table.Deck.AllColumns)

	_, err := stmt.Exec(h.app.DB)

	if err != nil {
		log.Println(err)
		return HTML(c, decks.CreateDeckForm(err))
	}

	c.Response().Header().Set("HX-Redirect", "/decks?toast_message=Deleted&toast_type=success")

	return c.NoContent(http.StatusCreated)
}
