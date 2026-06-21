package card

import (
	"context"

	"github.com/google/uuid"
	"github.com/valenoirs/flashcard-api/internal/domain"
)

type UpdateCardCommand struct {
	ID    uuid.UUID
	Front string
	Back  string
	Note  *string
	Class *string
}

func (u *UpdateCardCommand) ToDomain() *domain.Card {
	return &domain.Card{
		ID:    u.ID,
		Front: u.Front,
		Back:  u.Back,
		Note:  u.Note,
		Class: u.Class,
	}
}

type UpdateCardHandler struct {
	cardRepo domain.CardRepository
}

func NewUpdateCardHandler(
	cardRepo domain.CardRepository,
) *UpdateCardHandler {
	return &UpdateCardHandler{
		cardRepo: cardRepo,
	}
}

func (u *UpdateCardHandler) Handle(
	ctx context.Context,
	command *UpdateCardCommand,
) error {
	return u.cardRepo.UpdateCard(ctx, command.ToDomain())
}
