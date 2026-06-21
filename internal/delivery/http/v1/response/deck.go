package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/valenoirs/flashcard-api/internal/domain"
)

type GetDeckResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewGetDeckResponse(d *domain.Deck) GetDeckResponse {
	return GetDeckResponse{
		ID:        d.ID,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
