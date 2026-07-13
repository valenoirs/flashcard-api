package card

import (
	"context"

	"github.com/valenoirs/flashcard-api/internal/domain"
)


type UpdateCardHandler struct {
	cardRepo domain.CardRepository
}

func NewUpdateCardHandler(cardRepo domain.CardRepository) *UpdateCardHandler {
	return &UpdateCardHandler{
		cardRepo: cardRepo,
	}
}

func (h *UpdateCardHandler) Handle(ctx context.Context, cmd *UpdateCardCommand) error {
	card := &domain.Card{
		ID:      cmd.ID,
		Vocab:   cmd.Vocab,
		Kana:    cmd.Kana,
		Meaning: cmd.Meaning,
		English: cmd.English,
	}

	return h.cardRepo.UpdateCard(ctx, card)
}
