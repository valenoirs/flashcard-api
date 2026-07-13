package card

import "github.com/google/uuid"


type GetCardListQuery struct {
	DeckID uuid.UUID
}

type GetCardDetailQuery struct {
	ID uuid.UUID
}
