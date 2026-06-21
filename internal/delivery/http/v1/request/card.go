package request

import (
	"github.com/google/uuid"
	"github.com/valenoirs/flashcard-api/internal/usecase/card"
)

type CreateCardRequest struct {
	DeckID uuid.UUID `json:"deck_id" validate:"required"`
	Front  string    `json:"front" validate:"required"`
	Back   string    `json:"back" validate:"required"`
	Note   *string   `json:"note,omitempty"`
	Class  *string   `json:"class,omitempty"`
}

func NewCreateCardCommand(req *CreateCardRequest, cardID uuid.UUID) *card.CreateCardCommand {
	return &card.CreateCardCommand{
		ID: cardID,
		DeckID: req.DeckID,
		Front:  req.Front,
		Back:   req.Back,
		Note:   req.Note,
		Class:  req.Class,
	}
}

type UpdateCardRequest struct {
	Front string  `json:"front" validate:"required"`
	Back  string  `json:"back" validate:"required"`
	Note  *string `json:"note,omitempty"`
	Class *string `json:"class,omitempty"`
}

func NewUpdateCardCommand(req *UpdateCardRequest, cardID uuid.UUID) *card.UpdateCardCommand {
	return &card.UpdateCardCommand{
		ID:    cardID,
		Front: req.Front,
		Back:  req.Back,
		Note:  req.Note,
		Class: req.Class,
	}
}

func NewDeleteCardCommand(cardID uuid.UUID) *card.DeleteCardCommand {
	return &card.DeleteCardCommand{
		ID: cardID,
	}
}

func NewGetCardListQuery(deckID uuid.UUID) *card.GetCardListQuery {
	return &card.GetCardListQuery{
		DeckID: deckID,
	}
}

func NewGetCardDetailQuery(cardID uuid.UUID) *card.GetCardDetailQuery {
	return &card.GetCardDetailQuery{
		ID: cardID,
	}
}
