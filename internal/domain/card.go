package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID        uuid.UUID
	DeckID    uuid.UUID
	Front     string
	Back      string
	Note      *string
	Class     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CardRepository interface {
	CreateCard(ctx context.Context, card *Card) error
	UpdateCard(ctx context.Context, card *Card) error
	DeleteCard(ctx context.Context, card *Card) error
	GetCardList(ctx context.Context, id uuid.UUID) ([]Card, error)
	GetCardByID(ctx context.Context, id uuid.UUID) (*Card, error)
}
