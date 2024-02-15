package dao

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lewisd1996/baozi-zhongwen/internal/testutils"
)

func TestCardDao(t *testing.T) {
	ctx := context.Background()
	db, cleanup := testutils.InitializeTestDB(ctx, t)
	defer cleanup()
	dao := NewDao(db)

	t.Run("CardDao", func(t *testing.T) {
		var cardId uuid.UUID

		t.Run("Should create a card", func(t *testing.T) {
			card, err := dao.CreateCard("苹果", "Apple", testutils.TestDeckId)
			if err != nil {
				t.Fatalf("Failed to create card: %v", err)
			}
			cardId = card.ID
		})
		t.Run("Should get a card", func(t *testing.T) {
			card, err := dao.GetCardById(cardId.String(), testutils.TestDeckId)
			if err != nil {
				t.Fatalf("Failed to get card: %v", err)
			}
			if card.ID != cardId {
				t.Fatalf("Expected card ID to be %s, got %s", cardId, card.ID)
			}
			if card.Content != "苹果" {
				t.Fatalf("Expected card front to be 苹果, got %s", card.Content)
			}
			if card.Translation != "Apple" {
				t.Fatalf("Expected card back to be Apple, got %s", card.Translation)
			}
		})
		t.Run("Should update a card", func(t *testing.T) {
			card, err := dao.UpdateCard("橘子", "Orange", cardId.String(), testutils.TestDeckId)
			if err != nil {
				t.Fatalf("Failed to update card: %v", err)
			}
			if card.Content != "橘子" {
				t.Fatalf("Expected card front to be 橘子, got %s", card.Content)
			}
			if card.Translation != "Orange" {
				t.Fatalf("Expected card back to be Orange, got %s", card.Translation)
			}
		})
		t.Run("Should delete a card", func(t *testing.T) {
			err := dao.DeleteCard(cardId.String(), testutils.TestDeckId)
			if err != nil {
				t.Fatalf("Failed to delete card: %v", err)
			}
		})
	})
}
