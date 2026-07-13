package card

import (
	"context"

	"github.com/valenoirs/flashcard-api/internal/domain"
)


type DeleteCardHandler struct {
	cardRepo domain.CardRepository
}

func NewDeleteCardHandler(cardRepo domain.CardRepository) *DeleteCardHandler {
	return &DeleteCardHandler{
		cardRepo: cardRepo,
	}
}

func (c *DeleteCardHandler) Handle(ctx context.Context, command *DeleteCardCommand) error {
	return c.cardRepo.DeleteCard(ctx, command.ID)
}
