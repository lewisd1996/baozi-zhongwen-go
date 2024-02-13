//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
	"time"
)

type CardLearningProgress struct {
	ID             uuid.UUID `sql:"primary_key"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	SessionID      uuid.UUID
	CardID         uuid.UUID
	UserID         uuid.UUID
	ReviewCount    int32
	SuccessCount   int32
	LastReviewedAt time.Time
}