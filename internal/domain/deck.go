package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Deck struct {
	ID        uuid.UUID
	Name      string
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DeckRepository interface {
	CreateDeck(ctx context.Context, deck *Deck) error
	UpdateDeck(ctx context.Context, deck *Deck) error
	DeleteDeck(ctx context.Context, deckID uuid.UUID) error
	GetDeckList(ctx context.Context) ([]Deck, error)
}
