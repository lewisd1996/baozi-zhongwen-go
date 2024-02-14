package daos

import (
	"log"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	. "github.com/lewisd1996/baozi-zhongwen/sql/.jet/bz/public/model"
	"github.com/lewisd1996/baozi-zhongwen/sql/.jet/bz/public/table"
)

func (dao *Dao) GetCardById(cardId, deckId string) (Card, error) {
	stmt := table.Card.SELECT(table.Card.AllColumns).
		WHERE(
			table.Card.DeckID.EQ(UUID(uuid.MustParse(deckId))).
				AND(table.Card.ID.EQ(UUID(uuid.MustParse(cardId)))),
		).LIMIT(1)

	var card Card
	err := stmt.Query(dao.DB, &card)

	if err != nil {
		log.Println(err)
		return Card{}, err
	}

	return card, nil
}

func (dao *Dao) CreateCard(content, translation, deckId string) (Card, error) {
	card := Card{
		Content:     content,
		Translation: translation,
		DeckID:      uuid.MustParse(deckId),
	}

	stmt := table.Card.INSERT(table.Card.Content, table.Card.Translation, table.Card.DeckID).MODEL(card).RETURNING(table.Card.AllColumns)

	var cardRes Card
	err := stmt.Query(dao.DB, &cardRes)

	if err != nil {
		log.Println(err)
		return Card{}, err
	}

	return cardRes, nil
}

func (dao *Dao) UpdateCard(content, translation, cardId, deckId string) (Card, error) {
	updatedCard := Card{
		Content:     content,
		Translation: translation,
	}

	updateStmt := table.Card.UPDATE(table.Card.Content, table.Card.Translation).MODEL(updatedCard).WHERE(table.Card.ID.EQ(UUID(uuid.MustParse(cardId))).AND(table.Card.DeckID.EQ(UUID(uuid.MustParse(deckId))))).RETURNING(table.Card.AllColumns)

	var cardRes Card
	err := updateStmt.Query(dao.DB, &cardRes)

	if err != nil {
		return Card{}, err
	}

	return cardRes, nil
}

func (dao *Dao) GetCardsByDeckId(deckId string) ([]Card, error) {
	stmt := table.Card.SELECT(table.Card.AllColumns).WHERE(table.Card.DeckID.EQ(UUID(uuid.MustParse(deckId)))).ORDER_BY(table.Card.UpdatedAt.DESC())
	var cardsRes []Card
	err := stmt.Query(dao.DB, &cardsRes)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return cardsRes, nil
}

func (dao *Dao) GetUserCardCount(userId string) (int, error) {
	stmt := table.Card.SELECT(COUNT(table.Card.ID)).FROM(table.Card, table.Deck).WHERE(table.Card.DeckID.EQ(table.Deck.ID).AND(table.Deck.OwnerID.EQ(UUID(uuid.MustParse(userId))))).GROUP_BY(table.Card.DeckID).LIMIT(1)
	var res struct {
		Count int
	}
	err := stmt.Query(dao.DB, &res)

	if err != nil {
		// TODO: Handle this better
		if err.Error() == "qrm: no rows in result set" {
			return 0, nil
		}
		return 0, err
	}
	return res.Count, nil
}
