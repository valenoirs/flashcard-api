package deck

import (
	"context"

	"github.com/google/uuid"
	"github.com/valenoirs/flashcard-api/internal/domain"
)

type CreateDeckCommand struct {
	DeckID uuid.UUID
	Name   string
}

type CreateDeckHandler struct {
	deckRepo domain.DeckRepository
}

func NewCreateDeckHandler(deckRepo domain.DeckRepository) *CreateDeckHandler {
	return &CreateDeckHandler{
		deckRepo: deckRepo,
	}
}

func (c *CreateDeckHandler) Handle(ctx context.Context, cmd *CreateDeckCommand) error {
	deck := &domain.Deck{
		ID:   cmd.DeckID,
		Name: cmd.Name,
	}

	return c.deckRepo.CreateDeck(ctx, deck)
}
