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
	cardSentence := make([]*domain.Sentence, len(cmd.Sentences))
	for i := range cmd.Sentences {
		s := cmd.Sentences[i]
		cardSentence[i] = &domain.Sentence{
			CardID:   cmd.ID,
			Position: s.Position,
			Text:     s.Text,
			Reading:  s.Reading,
			IsTarget: s.IsTarget,
		}
	}

	card := &domain.Card{
		ID:       cmd.ID,
		DeckID:   cmd.DeckID,
		Vocab:    cmd.Vocab,
		Kana:     cmd.Kana,
		Meaning:  cmd.Meaning,
		English:  cmd.English,
		Sentences: cardSentence,
	}

	return h.cardRepo.CreateCard(ctx, card)
}
