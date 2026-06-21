package deck

import (
	"context"

	"github.com/valenoirs/flashcard-api/internal/domain"
)

type CreateDeckCommand struct {
	Name string
}

func (c *CreateDeckCommand) ToDomain() *domain.Deck {
	return &domain.Deck{
		Name: c.Name,
	}
}

type CreateDeckHandler struct {
	deckRepo domain.DeckRepository
}

func NewCreateDeckHandler(
	deckRepo domain.DeckRepository,
) *CreateDeckHandler {
	return &CreateDeckHandler{
		deckRepo: deckRepo,
	}
}

func (c *CreateDeckHandler) Handle(
	ctx context.Context,
	command *CreateDeckCommand,
) error {
	return c.deckRepo.CreateDeck(ctx, command.ToDomain())
}
