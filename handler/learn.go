package handler

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/app"
	"github.com/lewisd1996/baozi-zhongwen/view/learn"
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

	learningSession, err := h.app.Dao.GetLearningSessionById(deckId)

	if err != nil {
		// If no learning session exists, create one
		if err.Error() == "qrm: no rows in result set" {
			learningSession, err = h.app.Dao.CreateLearningSession(c.Request().Context(), deckId, userId)
		} else {
			return c.Redirect(302, "/decks")
		}
	}

	if learningSession.EndedAt != nil {
		println("Session ended")
		return c.Redirect(302, "/decks")
	}

	nextCard, err := h.app.Dao.GetNextLearningSessionCard(learningSession.ID.String(), userId)

	if err != nil {
		println("Error getting next card:", err.Error())
		return c.Redirect(302, "/decks")
	}

	// Get 3 other cards to be incorrect options
	incorrectOptions, err := h.app.Dao.GetLearningSessionIncorrectOptions(learningSession.ID.String(), userId, nextCard.Card.ID.String())

	var options []learn.LearnOption
	options = append(incorrectOptions, learn.LearnOption{Translation: nextCard.Translation, Correct: true})

	// Shuffle options
	for i := range options {
		j := i + rand.Intn(len(options)-i)
		options[i], options[j] = options[j], options[i]
	}

	state := learn.LearnState{
		SessionID:   learningSession.ID.String(),
		CardID:      nextCard.CardLearningProgress.CardID.String(),
		Content:     nextCard.Card.Content,
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

	learningSession, err := h.app.Dao.UpdateCardLearningProgress(c.Request().Context(), cardId, parsedCorrect, sessionId, userId)

	if err != nil {
		println("Error updating card learning progress:", err.Error())
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

	learningSession, err := h.app.Dao.GetLearningSessionById(sessionId)
	if err != nil {
		println("Error getting learning session:", err.Error())
		return c.Redirect(302, "/decks")
	}
	if learningSession.EndedAt != nil {
		println("Session ended")
		return c.Redirect(302, "/decks")
	}

	cards, err := h.app.Dao.GetLearningSessionCards(sessionId)
	if err != nil {
		println("Error getting learning session cards:", err.Error())
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
