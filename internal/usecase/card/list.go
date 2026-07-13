package card

import (
	"context"

	"github.com/valenoirs/flashcard-api/internal/domain"
)

type GetCardListHandler struct {
	cardRepo domain.CardRepository
}

func NewGetCardListHandler(cardRepo domain.CardRepository) *GetCardListHandler {
	return &GetCardListHandler{
		cardRepo: cardRepo,
	}
}

func (c *GetCardListHandler) Handle(ctx context.Context, q *GetCardListQuery) ([]domain.Card, error) {
	return c.cardRepo.GetCardList(ctx, q.DeckID)
}
