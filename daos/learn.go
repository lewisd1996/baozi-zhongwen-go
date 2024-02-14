package daos

import (
	"context"
	"fmt"
	"log"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	. "github.com/lewisd1996/baozi-zhongwen/sql/.jet/bz/public/model"
	"github.com/lewisd1996/baozi-zhongwen/sql/.jet/bz/public/table"
	"github.com/lewisd1996/baozi-zhongwen/view/learn"
)

// Get next card to learn
type CardWithLearningProgress struct {
	CardLearningProgress
	Card
}

func (dao *Dao) GetLearningSessionByDeckId(deckId, userId string) (LearningSession, error) {
	stmt := table.LearningSession.SELECT(table.LearningSession.AllColumns).WHERE(table.LearningSession.DeckID.EQ(UUID(uuid.MustParse(deckId))).AND(table.LearningSession.EndedAt.IS_NULL())).LIMIT(1)
	var sessionRes LearningSession
	err := stmt.Query(dao.DB, &sessionRes)

	if err != nil {
		log.Println(err)
		return LearningSession{}, err
	}

	return sessionRes, nil
}

func (dao *Dao) GetActiveLearningSessionById(sessionId string) (LearningSession, error) {
	stmt := table.LearningSession.SELECT(table.LearningSession.AllColumns).WHERE(table.LearningSession.ID.EQ(UUID(uuid.MustParse(sessionId))).AND(table.LearningSession.EndedAt.IS_NULL())).LIMIT(1)
	debugSql := stmt.DebugSql()
	log.Println(debugSql)
	var sessionRes LearningSession
	err := stmt.Query(dao.DB, &sessionRes)

	if err != nil {
		log.Println("Error getting learning session:", err.Error())
		return LearningSession{}, err
	}

	return sessionRes, nil
}

func (dao *Dao) GetEndedLearningSessionById(sessionId string) (LearningSession, error) {
	stmt := table.LearningSession.SELECT(table.LearningSession.AllColumns).WHERE(table.LearningSession.ID.EQ(UUID(uuid.MustParse(sessionId))).AND(table.LearningSession.EndedAt.IS_NOT_NULL())).LIMIT(1)
	debugSql := stmt.DebugSql()
	log.Println(debugSql)
	var sessionRes LearningSession
	err := stmt.Query(dao.DB, &sessionRes)

	if err != nil {
		log.Println("Error getting learning session:", err.Error())
		return LearningSession{}, err
	}

	return sessionRes, nil
}

func (dao *Dao) CreateLearningSession(ctx context.Context, deckId, userId string) (LearningSession, error) {
	var cards []Card
	var learningSession LearningSession

	// Begin transaction
	tx, err := dao.DB.Begin()
	// Rollback transaction if error
	defer tx.Rollback()

	if err != nil {
		log.Println(err)
		return LearningSession{}, err
	}

	stmt := table.LearningSession.INSERT(table.LearningSession.DeckID, table.LearningSession.UserID).VALUES(UUID(uuid.MustParse(deckId)), UUID(uuid.MustParse(userId))).RETURNING(table.LearningSession.AllColumns)

	err = stmt.QueryContext(ctx, tx, &learningSession)

	if err != nil {
		println("Error inserting learning session:", err.Error())
		err = tx.Rollback()
		return LearningSession{}, err
	}

	// Get 4 cards to start learning session
	newSessionCardsStmt := table.Card.SELECT(table.Card.AllColumns, table.CardLearningProgress.LastReviewedAt).FROM(table.Card.LEFT_JOIN(table.CardLearningProgress, table.CardLearningProgress.CardID.EQ(table.Card.ID))).WHERE(table.Card.DeckID.EQ(UUID(uuid.MustParse(deckId)))).ORDER_BY(table.CardLearningProgress.LastReviewedAt.DESC()).LIMIT(4)
	err = newSessionCardsStmt.QueryContext(ctx, tx, &cards)
	if err != nil {
		println("Error getting cards:", err.Error())
		err = tx.Rollback()
		return LearningSession{}, err
	}
	if len(cards) < 4 {
		println("Not enough cards to start learning session")
		err = tx.Rollback()
		return LearningSession{}, fmt.Errorf("Not enough cards to start learning session")
	}

	// Insert cards into learning progress
	for _, card := range cards {
		fmt.Println("Inserting card into session", card.ID, learningSession.ID)
		cardLearningProgressStmt := table.CardLearningProgress.INSERT(table.CardLearningProgress.SessionID, table.CardLearningProgress.CardID, table.CardLearningProgress.UserID).VALUES(learningSession.ID, card.ID, UUID(uuid.MustParse(userId)))
		_, err := cardLearningProgressStmt.ExecContext(ctx, tx)
		if err != nil {
			println("Error inserting card learning progress:", err.Error())
			err = tx.Rollback()
			return LearningSession{}, err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		println("Error committing transaction:", err.Error())
		return LearningSession{}, err
	}

	return learningSession, nil
}

func (dao *Dao) GetNextLearningSessionCard(sessionId, userId string) (CardWithLearningProgress, error) {
	// Get next card to learn
	var nextCardDest CardWithLearningProgress
	nextCardStmt := table.CardLearningProgress.SELECT(table.CardLearningProgress.AllColumns, table.Card.Content, table.Card.Translation).FROM(table.CardLearningProgress.INNER_JOIN(table.Card, table.Card.ID.EQ(table.CardLearningProgress.CardID))).WHERE(table.CardLearningProgress.SessionID.EQ(UUID(uuid.MustParse(sessionId))).AND(table.CardLearningProgress.UserID.EQ(UUID(uuid.MustParse(userId))))).ORDER_BY(table.CardLearningProgress.LastReviewedAt.ASC()).LIMIT(1)
	err := nextCardStmt.Query(dao.DB, &nextCardDest)
	if err != nil {
		log.Println(err)
		return nextCardDest, err
	}

	return nextCardDest, nil
}

func (dao *Dao) GetLearningSessionIncorrectOptions(nextCardId, sessionId, userId uuid.UUID) ([]learn.LearnOption, error) {
	// Get 3 other cards to be incorrect options
	var incorrectCardsDest []CardWithLearningProgress
	incorrectCardsStmt := table.CardLearningProgress.SELECT(table.CardLearningProgress.AllColumns, table.Card.Content, table.Card.Translation).FROM(table.CardLearningProgress.INNER_JOIN(table.Card, table.Card.ID.EQ(table.CardLearningProgress.CardID))).WHERE(table.CardLearningProgress.SessionID.EQ(UUID(sessionId)).AND(table.CardLearningProgress.UserID.EQ(UUID(userId))).AND(table.CardLearningProgress.CardID.NOT_EQ(UUID(nextCardId)))).LIMIT(3)
	err := incorrectCardsStmt.Query(dao.DB, &incorrectCardsDest)
	if err != nil {
		log.Println("Error getting incorrect cards:", err.Error())
		return nil, err
	}

	var options []learn.LearnOption
	for _, card := range incorrectCardsDest {
		options = append(options, learn.LearnOption{Translation: card.Card.Translation, Correct: false})
	}

	return options, nil
}

func (dao *Dao) GetCardLearningProgress(cardId, sessionId, userId string) (CardLearningProgress, error) {
	var cardLearningProgress CardLearningProgress
	selectStmt := table.CardLearningProgress.SELECT(table.CardLearningProgress.AllColumns).WHERE(table.CardLearningProgress.SessionID.EQ(UUID(uuid.MustParse(sessionId))).AND(table.CardLearningProgress.CardID.EQ(UUID(uuid.MustParse(cardId)))).AND(table.CardLearningProgress.UserID.EQ(UUID(uuid.MustParse(userId))))).LIMIT(1)
	err := selectStmt.Query(dao.DB, &cardLearningProgress)

	if err != nil {
		println("Error getting card learning progress:", err.Error())
		return CardLearningProgress{}, err
	}

	return cardLearningProgress, nil
}

// TODO: Refactor, too many responsibilities
func (dao *Dao) UpdateCardLearningProgress(ctx context.Context, cardId string, isCorrect bool, sessionId, userId string) (LearningSession, error) {
	// Get learning session
	learningSession, err := dao.GetActiveLearningSessionById(sessionId)

	if err != nil {
		println("Error getting learning session:", err.Error())
		return LearningSession{}, err
	}
	if learningSession.EndedAt != nil {
		println("Learning session has ended")
		return LearningSession{}, fmt.Errorf("Learning session has ended")
	}

	// Get card learning progress
	cardLearningProgress, err := dao.GetCardLearningProgress(cardId, sessionId, userId)
	if err != nil {
		println("Error getting card learning progress:", err.Error())
		return LearningSession{}, err
	}

	// Begin transaction
	tx, err := dao.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		println("Error starting transaction:", err.Error())
		return LearningSession{}, err
	}

	// Update session review count
	updateReviewCountStmt := table.LearningSession.UPDATE(table.LearningSession.ReviewCount).SET(table.LearningSession.ReviewCount.ADD(Int32(1))).WHERE(table.LearningSession.ID.EQ(UUID(cardLearningProgress.SessionID))).RETURNING(table.LearningSession.AllColumns)
	err = updateReviewCountStmt.QueryContext(ctx, tx, &learningSession)
	if err != nil {
		println("Updating this session:", learningSession.ID.String())
		println("Error updating review count:", err.Error())
		return LearningSession{}, err
	}

	// Update card learning progress
	cardLearningProgress.LastReviewedAt = time.Now()
	cardLearningProgress.ReviewCount = cardLearningProgress.ReviewCount + 1
	if isCorrect {
		cardLearningProgress.SuccessCount = cardLearningProgress.SuccessCount + 1
	}
	updatedCardLearningProgressFields := CardLearningProgress{
		LastReviewedAt: time.Now(),
		ReviewCount:    cardLearningProgress.ReviewCount,
		SuccessCount:   cardLearningProgress.SuccessCount,
	}
	updateCardLearningProgressStmt := table.CardLearningProgress.UPDATE(table.CardLearningProgress.LastReviewedAt, table.CardLearningProgress.ReviewCount, table.CardLearningProgress.SuccessCount).MODEL(updatedCardLearningProgressFields).WHERE(table.CardLearningProgress.SessionID.EQ(UUID(cardLearningProgress.SessionID)).AND(table.CardLearningProgress.CardID.EQ(UUID(cardLearningProgress.CardID)))).RETURNING(table.CardLearningProgress.AllColumns)
	err = updateCardLearningProgressStmt.QueryContext(ctx, tx, &cardLearningProgress)

	// Update session progress
	if learningSession.ReviewCount == 10 {
		endedAt := time.Now()
		updatedLearningSessionFields := LearningSession{
			EndedAt: &endedAt,
		}
		updateLearningSessionStmt := table.LearningSession.UPDATE(table.LearningSession.EndedAt).MODEL(updatedLearningSessionFields).WHERE(table.LearningSession.ID.EQ(UUID(learningSession.ID))).RETURNING(table.LearningSession.AllColumns)
		err = updateLearningSessionStmt.QueryContext(ctx, tx, &learningSession)
		if err != nil {
			println("Error ending learning session:", err.Error())
			return LearningSession{}, err
		}
	}

	// Commit transaction
	err = tx.Commit()

	return learningSession, nil
}

func (dao *Dao) GetLearningSessionCards(sessionId string) ([]learn.SessionCard, error) {
	var cards []learn.SessionCard
	selectStmt := table.CardLearningProgress.SELECT(table.CardLearningProgress.AllColumns, table.Card.Content, table.Card.Translation).FROM(table.CardLearningProgress.INNER_JOIN(table.Card, table.Card.ID.EQ(table.CardLearningProgress.CardID))).WHERE(table.CardLearningProgress.SessionID.EQ(UUID(uuid.MustParse(sessionId)))).ORDER_BY(table.CardLearningProgress.LastReviewedAt.DESC())
	err := selectStmt.Query(dao.DB, &cards)
	if err != nil {
		println("Error getting cards:", err.Error())
		return nil, err
	}

	return cards, nil
}

func (dao *Dao) GetUserCompletedLearningSessionCount(userId string) (int, error) {
	stmt := table.LearningSession.SELECT(COUNT(table.LearningSession.ID)).FROM(table.LearningSession).WHERE(table.LearningSession.UserID.EQ(UUID(uuid.MustParse(userId))).AND(table.LearningSession.EndedAt.IS_NOT_NULL())).LIMIT(1)
	var res struct {
		Count int
	}
	err := stmt.Query(dao.DB, &res)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return res.Count, nil
}
