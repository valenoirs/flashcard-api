package card

import (
	"context"

	"github.com/valenoirs/flashcard-api/internal/domain"
)

type CreateCardHandler struct {
	cardRepo domain.CardRepository
}

func NewCreateCardHandler(cardRepo domain.CardRepository) *CreateCardHandler {
	return &CreateCardHandler{
		cardRepo: cardRepo,
	}
}

func (h *CreateCardHandler) Handle(ctx context.Context, cmd *CreateCardCommand) error {
	card := &domain.Card{
		ID:       cmd.ID,
		DeckID:   cmd.DeckID,
		Vocab:    cmd.Vocab,
		Kana:     cmd.Kana,
		Meaning:  cmd.Meaning,
		English:  cmd.English,
		Sentence: cmd.Sentence,
	}

	return h.cardRepo.CreateCard(ctx, card)
}
