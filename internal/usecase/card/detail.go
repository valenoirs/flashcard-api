package card

import (
	"context"

	"github.com/google/uuid"
	"github.com/valenoirs/flashcard-api/internal/domain"
)

type GetCardDetailQuery struct {
	ID uuid.UUID
}

type GetCardDetailHandler struct {
	cardRepo domain.CardRepository
}	

func NewGetCardDetailHandler(
	cardRepo domain.CardRepository,
) *GetCardDetailHandler {
	return &GetCardDetailHandler{
		cardRepo: cardRepo,
	}
}

func (c *GetCardDetailHandler) Handle(ctx context.Context, q *GetCardDetailQuery) (*domain.Card, error){
	return c.cardRepo.GetCardByID(ctx, q.ID)
}
