package deck

import (
	"context"

	"github.com/valenoirs/flashcard-api/internal/domain"
)

type GetDeckListQuery struct{}

func (d *GetDeckListQuery) ToDomain() *domain.Deck {
	return &domain.Deck{}
}

type GetDeckListHandler struct {
	deckRepo domain.DeckRepository
}

func NewGetDeckListHandler(deckRepo domain.DeckRepository) *GetDeckListHandler {
	return &GetDeckListHandler{
		deckRepo: deckRepo,
	}
}

func (d *GetDeckListHandler) Handle(ctx context.Context, query *GetDeckListQuery) ([]domain.Deck, error) {
	return d.deckRepo.GetDeckList(ctx)
}
