package card

import (
	"context"

	"github.com/google/uuid"
	"github.com/valenoirs/flashcard-api/internal/domain"
)

type CreateCardCommand struct {
	ID     uuid.UUID
	DeckID uuid.UUID
	Front  string
	Back   string
	Note   *string
	Class  *string
}

func (c *CreateCardCommand) ToDomain() *domain.Card {
	return &domain.Card{
		ID:     c.ID,
		DeckID: c.DeckID,
		Front:  c.Front,
		Back:   c.Back,
		Note:   c.Note,
		Class:  c.Class,
	}
}

type CreateCardHandler struct {
	cardRepo domain.CardRepository
}

func NewCreateCardHandler(
	cardRepo domain.CardRepository,
) *CreateCardHandler {
	return &CreateCardHandler{
		cardRepo: cardRepo,
	}
}

func (c *CreateCardHandler) Handle(
	ctx context.Context,
	command *CreateCardCommand,
) error {
	return c.cardRepo.CreateCard(ctx, command.ToDomain())
}
