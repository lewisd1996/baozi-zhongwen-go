package dao

import (
	"log"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	. "github.com/lewisd1996/baozi-zhongwen/sql/.jet/bz/public/model"
	"github.com/lewisd1996/baozi-zhongwen/sql/.jet/bz/public/table"
)

func (dao *Dao) GetDeckById(deckId string) (Deck, error) {
	stmt := table.Deck.SELECT(table.Deck.AllColumns).WHERE(table.Deck.ID.EQ(UUID(uuid.MustParse(deckId)))).LIMIT(1)
	var deckRes Deck
	err := stmt.Query(dao.DB, &deckRes)

	if err != nil {
		log.Println(err)
		return Deck{}, err
	}

	return deckRes, nil
}

func (dao *Dao) GetDecksByOwnerId(ownerId string) ([]Deck, error) {
	stmt := table.Deck.SELECT(table.Deck.AllColumns).WHERE(table.Deck.OwnerID.EQ(UUID(uuid.MustParse(ownerId)))).ORDER_BY(table.Deck.UpdatedAt)
	var decksRes []Deck
	err := stmt.Query(dao.DB, &decksRes)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return decksRes, nil
}

func (dao *Dao) CreateDeck(name, description, ownerId string) (Deck, error) {
	deck := Deck{
		Description: description,
		Name:        name,
		OwnerID:     uuid.MustParse(ownerId),
	}
	stmt := table.Deck.INSERT(table.Deck.OwnerID, table.Deck.Name, table.Deck.Description).MODEL(deck).RETURNING(table.Deck.AllColumns)
	var deckRes Deck
	err := stmt.Query(dao.DB, &deckRes)

	if err != nil {
		log.Println(err)
		return Deck{}, err
	}

	return deckRes, nil
}

func (dao *Dao) DeleteDeck(deckId, ownerId string) error {
	stmt := table.Deck.DELETE().WHERE(table.Deck.OwnerID.EQ(UUID(uuid.MustParse(ownerId))).AND(table.Deck.ID.EQ(UUID(uuid.MustParse(deckId))))).RETURNING(table.Deck.AllColumns)
	_, err := stmt.Exec(dao.DB)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (dao *Dao) GetUserDeckCount(userId string) (int, error) {
	stmt := table.Card.SELECT(COUNT(table.Deck.ID)).FROM(table.Deck).WHERE(table.Deck.OwnerID.EQ(UUID(uuid.MustParse(userId)))).LIMIT(1)
	var res struct {
		Count int
	}
	err := stmt.Query(dao.DB, &res)
	if err != nil {
		log.Println("Error getting user deck count:", err)
		return 0, err
	}
	return res.Count, nil
}
