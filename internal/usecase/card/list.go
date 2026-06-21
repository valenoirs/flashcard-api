package card

import (
	"context"

	"github.com/google/uuid"
	"github.com/valenoirs/flashcard-api/internal/domain"
)

type GetCardListQuery struct {
	DeckID uuid.UUID
}

func (c *GetCardListQuery) ToDomain() *domain.Card {
	return &domain.Card{
		DeckID: c.DeckID,
	}
}

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
