package dao

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lewisd1996/baozi-zhongwen/internal/testutils"
)

func TestDeckDao(t *testing.T) {
	ctx := context.Background()
	db, cleanup := testutils.InitializeTestDB(ctx, t)
	defer cleanup()
	dao := NewDao(db)

	t.Run("DeckDao", func(t *testing.T) {
		var deckId uuid.UUID

		t.Run("Should create a deck", func(t *testing.T) {
			deck, err := dao.CreateDeck("Business Chinese", "A deck for learning business Chinese", testutils.TestUserId)
			if err != nil {
				t.Fatalf("Failed to create deck: %v", err)
			}
			deckId = deck.ID
		})
		t.Run("Should delete a deck", func(t *testing.T) {
			err := dao.DeleteDeck(deckId.String(), testutils.TestUserId)
			if err != nil {
				t.Fatalf("Failed to delete deck: %v", err)
			}
		})
	})
}
