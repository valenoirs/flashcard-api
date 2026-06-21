package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/valenoirs/flashcard-api/internal/domain"
)

type GetCardResponse struct {
	ID        uuid.UUID `json:"id"`
	Front     string    `json:"front"`
	Back      string    `json:"back"`
	Note      *string   `json:"note"`
	Class     *string   `json:"class"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewGetCardResponse(c *domain.Card) GetCardResponse {
	return GetCardResponse{
		ID:        c.ID,
		Front:     c.Front,
		Back:      c.Back,
		Note:      c.Note,
		Class:     c.Class,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
