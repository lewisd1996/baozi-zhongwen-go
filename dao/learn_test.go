package dao

import (
	"context"
	"slices"
	"testing"

	"github.com/lewisd1996/baozi-zhongwen/internal/testutils"
	"github.com/lewisd1996/baozi-zhongwen/sql/.jet/bz/public/model"
)

func TestLearnDao(t *testing.T) {
	ctx := context.Background()
	db, cleanup := testutils.InitializeTestDB(ctx, t)
	defer cleanup()
	dao := NewDao(db)

	t.Run("LearnDao", func(t *testing.T) {
		var learningSessionId string
		var card1, card2, card3, card4 model.Card

		t.Run("Should fail to create a learning session when there are < 4 cards in a deck", func(t *testing.T) {
			_, err := dao.CreateLearningSession(ctx, testutils.TestDeckId, testutils.TestUserId)
			if err == nil {
				t.Fatalf("Created learning session without 4 cards")
			}
		})
		t.Run("Should create a learning session", func(t *testing.T) {
			// Create a deck with 4 cards
			deck, err := dao.CreateDeck("Test Deck", "Test Deck Description", testutils.TestUserId)
			if err != nil {
				t.Fatalf("Failed to create deck: %v", err)
			}
			card1, err = dao.CreateCard("苹果", "Apple", deck.ID.String())
			if err != nil {
				t.Fatalf("Failed to create card: %v", err)
			}
			card2, err = dao.CreateCard("橘子", "Orange", deck.ID.String())
			if err != nil {
				t.Fatalf("Failed to create card: %v", err)
			}
			card3, err = dao.CreateCard("香蕉", "Banana", deck.ID.String())
			if err != nil {
				t.Fatalf("Failed to create card: %v", err)
			}
			card4, err = dao.CreateCard("梨", "Pear", deck.ID.String())
			if err != nil {
				t.Fatalf("Failed to create card: %v", err)
			}

			// Create a learning session
			learningSession, err := dao.CreateLearningSession(ctx, deck.ID.String(), testutils.TestUserId)
			if err != nil {
				t.Fatalf("Failed to create learning session: %v", err)
			}
			learningSessionId = learningSession.ID.String()
		})
		t.Run("Should get an active learning session", func(t *testing.T) {
			learningSession, err := dao.GetActiveLearningSessionById(learningSessionId)
			if err != nil {
				t.Fatalf("Failed to get learning session: %v", err)
			}
			if learningSession.ID.String() != learningSessionId {
				t.Fatalf("Expected learning session ID to be %s, got %s", learningSessionId, learningSession.ID.String())
			}
		})
		t.Run("Should get a learning session's next card", func(t *testing.T) {
			nextCard, err := dao.GetNextLearningSessionCard(learningSessionId, testutils.TestUserId)
			if err != nil {
				t.Fatalf("Failed to get next card: %v", err)
			}
			cardIds := []string{card1.ID.String(), card2.ID.String(), card3.ID.String(), card4.ID.String()}

			if !slices.Contains(cardIds, nextCard.CardID.String()) {
				t.Fatalf("Expected card ID to be one of %v, got %s", cardIds, nextCard.CardID.String())
			}
		})

	})
}
