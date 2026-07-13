package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDeckRequest struct {
	Name string `json:"name" validate:"required"`
}

type DeckResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
