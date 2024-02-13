package handler

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/app"
	"github.com/lewisd1996/baozi-zhongwen/view/learn"

	// dot import so that jet go code would resemble as much as native SQL
	// dot import is not mandatory

	. "github.com/go-jet/jet/v2/postgres"
	. "github.com/lewisd1996/baozi-zhongwen/sql/.jet/bz/public/model"
	"github.com/lewisd1996/baozi-zhongwen/sql/.jet/bz/public/table"
)

type LearnHandler struct {
	app *app.App
}

func NewLearnHandler(a *app.App) LearnHandler {
	return LearnHandler{app: a}
}

func (h LearnHandler) HandleLearnShow(c echo.Context) error {
	userId := c.Get("user_id").(string)
	deckId := c.QueryParam("deck_id")

	if deckId == "" {
		return c.Redirect(302, "/decks")
	}

	var learningSession LearningSession
	stmt := table.LearningSession.SELECT(table.LearningSession.AllColumns).WHERE(table.LearningSession.DeckID.EQ(UUID(uuid.MustParse(deckId))).AND(table.LearningSession.EndedAt.IS_NULL())).LIMIT(1)
	err := stmt.Query(h.app.DB, &learningSession)

	if err != nil {
		// If no learning session exists, create one
		if err.Error() == "qrm: no rows in result set" {
			var cards []Card

			// Begin transaction
			tx, err := h.app.DB.Begin()
			// Rollback transaction if error
			defer tx.Rollback()

			if err != nil {
				println("Error starting transaction:", err.Error())
				return c.Redirect(302, "/decks")
			}

			stmt := table.LearningSession.INSERT(table.LearningSession.DeckID, table.LearningSession.UserID).VALUES(UUID(uuid.MustParse(deckId)), UUID(uuid.MustParse(userId))).RETURNING(table.LearningSession.AllColumns)
			err = stmt.QueryContext(c.Request().Context(), tx, &learningSession)
			if err != nil {
				println("Error inserting learning session:", err.Error())
				err = tx.Rollback()
				return c.Redirect(302, "/decks")
			}

			// Get 4 cards to start learning session
			newSessionCardsStmt := table.Card.SELECT(table.Card.AllColumns, table.CardLearningProgress.LastReviewedAt).FROM(table.Card.LEFT_JOIN(table.CardLearningProgress, table.CardLearningProgress.CardID.EQ(table.Card.ID))).WHERE(table.Card.DeckID.EQ(UUID(uuid.MustParse(deckId)))).ORDER_BY(table.CardLearningProgress.LastReviewedAt.DESC()).LIMIT(4)
			debugQuery := newSessionCardsStmt.DebugSql()
			println("Debug query:", debugQuery)
			err = newSessionCardsStmt.QueryContext(c.Request().Context(), tx, &cards)
			if err != nil {
				println("Error getting cards:", err.Error())
				err = tx.Rollback()
				return c.Redirect(302, "/decks")
			}
			if len(cards) < 4 {
				println("Not enough cards to start learning session")
				err = tx.Rollback()
				return c.Redirect(302, "/decks")
			}

			// Insert cards into learning progress
			for _, card := range cards {
				fmt.Println("Inserting card into session", card.ID, learningSession.ID)
				cardLearningProgressStmt := table.CardLearningProgress.INSERT(table.CardLearningProgress.SessionID, table.CardLearningProgress.CardID, table.CardLearningProgress.UserID).VALUES(learningSession.ID, card.ID, UUID(uuid.MustParse(userId)))
				_, err := cardLearningProgressStmt.ExecContext(c.Request().Context(), tx)
				if err != nil {
					println("Error inserting card learning progress:", err.Error())
					err = tx.Rollback()
					return c.Redirect(302, "/decks")
				}
			}

			// Commit transaction
			err = tx.Commit()
			if err != nil {
				println("Error committing transaction:", err.Error())
				return c.Redirect(302, "/decks")
			}
		} else {
			return c.Redirect(302, "/decks")
		}
	}

	if learningSession.EndedAt != nil {
		println("Session ended")
		return c.Redirect(302, "/decks")
	}

	// Get next card to learn
	var nextCardDest struct {
		CardLearningProgress
		Card
	}
	nextCardStmt := table.CardLearningProgress.SELECT(table.CardLearningProgress.AllColumns, table.Card.Content, table.Card.Translation).FROM(table.CardLearningProgress.INNER_JOIN(table.Card, table.Card.ID.EQ(table.CardLearningProgress.CardID))).WHERE(table.CardLearningProgress.SessionID.EQ(UUID(uuid.MustParse(learningSession.ID.String()))).AND(table.CardLearningProgress.UserID.EQ(UUID(uuid.MustParse(userId))))).ORDER_BY(table.CardLearningProgress.LastReviewedAt.ASC()).LIMIT(1)
	err = nextCardStmt.Query(h.app.DB, &nextCardDest)
	if err != nil {
		println("Error in nextCardStmt:", err.Error())
		return c.Redirect(302, "/decks")
	}

	println("Next card:", nextCardDest.Card.Content, nextCardDest.Card.Translation)

	// Get 3 other cards to be incorrect options
	var incorrectCardsDest []struct {
		CardLearningProgress
		Card
	}
	incorrectCardsStmt := table.CardLearningProgress.SELECT(table.CardLearningProgress.AllColumns, table.Card.Content, table.Card.Translation).FROM(table.CardLearningProgress.INNER_JOIN(table.Card, table.Card.ID.EQ(table.CardLearningProgress.CardID))).WHERE(table.CardLearningProgress.SessionID.EQ(UUID(uuid.MustParse(learningSession.ID.String()))).AND(table.CardLearningProgress.UserID.EQ(UUID(uuid.MustParse(userId)))).AND(table.CardLearningProgress.CardID.NOT_EQ(UUID(nextCardDest.CardLearningProgress.CardID)))).LIMIT(3)
	err = incorrectCardsStmt.Query(h.app.DB, &incorrectCardsDest)
	if err != nil {
		println("Error in incorrectCardsStmt:", err.Error())
		return c.Redirect(302, "/decks")
	}
	println("Incorrect cards count:", len(incorrectCardsDest))

	var options []learn.LearnOption
	for _, card := range incorrectCardsDest {
		fmt.Println("Incorrect card:", card.Card.Content, card.Card.Translation)
		options = append(options, learn.LearnOption{Translation: card.Card.Translation, Correct: false})
	}
	options = append(options, learn.LearnOption{Translation: nextCardDest.Card.Translation, Correct: true})

	// Shuffle options
	for i := range options {
		j := i + rand.Intn(len(options)-i)
		options[i], options[j] = options[j], options[i]
	}

	state := learn.LearnState{
		SessionID:   learningSession.ID.String(),
		CardID:      nextCardDest.CardLearningProgress.CardID.String(),
		Content:     nextCardDest.Card.Content,
		Options:     options,
		ReviewCount: int(learningSession.ReviewCount),
	}
	return Render(c, learn.Show(userId, c.Path(), state))
}

func (h LearnHandler) HandleLearnSessionAnswerSubmit(c echo.Context) error {
	userId := c.Get("user_id").(string)
	sessionId := c.Param("learning_session_id")
	cardId, correct := c.QueryParam("card_id"), c.QueryParam("correct")

	if sessionId == "" || cardId == "" || correct == "" {
		println("Missing params")
		return c.Redirect(302, "/decks")
	}

	parsedCorrect, err := strconv.ParseBool(correct)
	if err != nil {
		println("Error parsing correct:", err.Error())
		return c.Redirect(302, "/decks")
	}

	// Get learning session
	var learningSession LearningSession
	selectStmt := table.LearningSession.SELECT(table.LearningSession.AllColumns).WHERE(table.LearningSession.ID.EQ(UUID(uuid.MustParse(sessionId)))).LIMIT(1)
	err = selectStmt.Query(h.app.DB, &learningSession)
	if err != nil || learningSession.EndedAt != nil {
		println("Error getting learning session:", err.Error())
		return c.Redirect(302, "/decks")
	}

	// Get card learning progress
	var cardLearningProgress CardLearningProgress
	selectStmt = table.CardLearningProgress.SELECT(table.CardLearningProgress.AllColumns).WHERE(table.CardLearningProgress.SessionID.EQ(UUID(uuid.MustParse(sessionId))).AND(table.CardLearningProgress.CardID.EQ(UUID(uuid.MustParse(cardId)))).AND(table.CardLearningProgress.UserID.EQ(UUID(uuid.MustParse(userId))))).LIMIT(1)
	err = selectStmt.Query(h.app.DB, &cardLearningProgress)
	if err != nil {
		println("Error getting card learning progress:", err.Error())
		return c.Redirect(302, "/decks")
	}

	// Start transaction
	tx, err := h.app.DB.Begin()
	if err != nil {
		println("Error starting transaction:", err.Error())
		return c.Redirect(302, "/decks")
	}

	// Update session review count
	updateReviewCountStmt := table.LearningSession.UPDATE(table.LearningSession.ReviewCount).SET(table.LearningSession.ReviewCount.ADD(Int32(1))).WHERE(table.LearningSession.ID.EQ(UUID(cardLearningProgress.SessionID))).RETURNING(table.LearningSession.AllColumns)
	err = updateReviewCountStmt.QueryContext(c.Request().Context(), tx, &learningSession)

	if err != nil {
		println("Updating this session:", learningSession.ID.String())
		println("Error updating review count:", err.Error())
		return c.Redirect(302, "/decks")
	}

	// Update card learning progress
	cardLearningProgress.LastReviewedAt = time.Now()
	cardLearningProgress.ReviewCount = cardLearningProgress.ReviewCount + 1
	if parsedCorrect {
		println("Correct!")
		cardLearningProgress.SuccessCount = cardLearningProgress.SuccessCount + 1
	}
	updatedCardLearningProgressFields := CardLearningProgress{
		LastReviewedAt: time.Now(),
		ReviewCount:    cardLearningProgress.ReviewCount,
		SuccessCount:   cardLearningProgress.SuccessCount,
	}
	updateCardLearningProgressStmt := table.CardLearningProgress.UPDATE(table.CardLearningProgress.LastReviewedAt, table.CardLearningProgress.ReviewCount, table.CardLearningProgress.SuccessCount).MODEL(updatedCardLearningProgressFields).WHERE(table.CardLearningProgress.SessionID.EQ(UUID(cardLearningProgress.SessionID)).AND(table.CardLearningProgress.CardID.EQ(UUID(cardLearningProgress.CardID)))).RETURNING(table.CardLearningProgress.AllColumns)
	err = updateCardLearningProgressStmt.QueryContext(c.Request().Context(), tx, &cardLearningProgress)

	if err != nil {
		println("Updating this card:", cardLearningProgress.CardID.String())
		println("Error updating card learning progress:", err.Error())
		return c.Redirect(302, "/decks")
	}

	// Update session progress
	if learningSession.ReviewCount == 10 {
		endedAt := time.Now()
		updatedLearningSessionFields := LearningSession{
			EndedAt: &endedAt,
		}
		updateLearningSessionStmt := table.LearningSession.UPDATE(table.LearningSession.EndedAt).MODEL(updatedLearningSessionFields).WHERE(table.LearningSession.ID.EQ(UUID(learningSession.ID))).RETURNING(table.LearningSession.AllColumns)
		err = updateLearningSessionStmt.QueryContext(c.Request().Context(), tx, &learningSession)
		if err != nil {
			println("Error ending learning session:", err.Error())
			return c.Redirect(302, "/decks")
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		println("Error committing transaction:", err.Error())
		return c.Redirect(302, "/decks")
	}

	if learningSession.ReviewCount == 10 {
		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/learn/%s/summary", learningSession.ID.String()))
	} else {
		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/learn?deck_id=%s", learningSession.DeckID.String()))
	}

	return c.NoContent(200)
}

func (h LearnHandler) HandleLearnSessionSummaryShow(c echo.Context) error {
	userId := c.Get("user_id").(string)

	sessionId := c.Param("learning_session_id")
	if sessionId == "" {
		return c.Redirect(302, "/decks")
	}

	var learningSession LearningSession
	selectStmt := table.LearningSession.SELECT(table.LearningSession.AllColumns).WHERE(table.LearningSession.ID.EQ(UUID(uuid.MustParse(sessionId)))).LIMIT(1)
	err := selectStmt.Query(h.app.DB, &learningSession)
	if err != nil {
		println("Error getting learning session:", err.Error())
		return c.Redirect(302, "/decks")
	}

	if learningSession.EndedAt == nil {
		return c.Redirect(302, "/decks")
	}

	var cards []learn.SessionCard
	selectStmt = table.CardLearningProgress.SELECT(table.CardLearningProgress.AllColumns, table.Card.Content, table.Card.Translation).FROM(table.CardLearningProgress.INNER_JOIN(table.Card, table.Card.ID.EQ(table.CardLearningProgress.CardID))).WHERE(table.CardLearningProgress.SessionID.EQ(UUID(uuid.MustParse(sessionId)))).ORDER_BY(table.CardLearningProgress.LastReviewedAt.DESC())
	err = selectStmt.Query(h.app.DB, &cards)
	if err != nil {
		println("Error getting cards:", err.Error())
		return c.Redirect(302, "/decks")
	}

	state := learn.SummaryState{
		DeckID:       learningSession.DeckID.String(),
		SessionID:    learningSession.ID.String(),
		SessionCards: cards,
		EndedAt:      learningSession.EndedAt.Format("2006-01-02"),
	}

	return Render(c, learn.ShowSummary(state, userId, c.Path()))
}
