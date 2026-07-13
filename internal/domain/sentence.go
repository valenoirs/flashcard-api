package domain

import (
	"context"

	"github.com/google/uuid"
)

type Sentence struct {
	ID       uuid.UUID
	CardID   uuid.UUID
	Position int
	Text     string
	Reading  string
	IsTarget bool
}

type SentenceRepository interface {
	SaveBatch(ctx context.Context, sentences []Sentence) error
}
