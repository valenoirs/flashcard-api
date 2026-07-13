package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request

type CreateCardRequest struct {
	DeckID    uuid.UUID              `json:"deck_id" validate:"required"`
	Vocab     string                 `json:"vocab" validate:"required"`
	Kana      string                 `json:"kana" validate:"required"`
	Meaning   string                 `json:"meaning" validate:"required"`
	English   string                 `json:"english" validate:"required"`
	Sentences []*CardSentenceRequest `json:"sentences" validate:"required"`
}

type CardSentenceRequest struct {
	Position int    `json:"position" validate:"required"`
	Text     string `json:"text" validate:"required"`
	Reading  string `json:"reading" validate:"required"`
	IsTarget bool   `json:"is_target" validate:"required"`
}

type UpdateCardRequest struct {
	Vocab   string `json:"vocab" validate:"required"`
	Kana    string `json:"kana" validate:"required"`
	Meaning string `json:"meaning" validate:"required"`
	English string `json:"english" validate:"required"`
}

// Response

type CardResponse struct {
	ID        uuid.UUID               `json:"id"`
	Vocab     string                  `json:"vocab"`
	Kana      string                  `json:"kana"`
	English   string                  `json:"english"`
	Meaning   string                  `json:"meaning"`
	Sentences []*CardSentenceResponse `json:"sentences"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
}

type CardSentenceResponse struct {
	ID       uuid.UUID `json:"id"`
	Position int       `json:"position"`
	Text     string    `json:"text"`
	Reading  string    `json:"reading"`
	IsTarget bool      `json:"is_target"`
}
