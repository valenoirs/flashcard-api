package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID        uuid.UUID
	DeckID    uuid.UUID
	Vocab     string
	Kana      string
	Meaning   string
	English   string
	Sentences  []*Sentence
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CardRepository interface {
	CreateCard(ctx context.Context, card *Card) error
	UpdateCard(ctx context.Context, card *Card) error
	DeleteCard(ctx context.Context, cardID uuid.UUID) error
	GetCardList(ctx context.Context, deckID uuid.UUID) ([]Card, error)
	GetCardByID(ctx context.Context, cardID uuid.UUID) (*Card, error)
}
