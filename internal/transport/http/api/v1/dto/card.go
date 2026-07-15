package dto

import "github.com/google/uuid"

// Request

type CreateCardRequest struct {
	DeckID      uuid.UUID `json:"deck_id" validate:"required"`
	Vocab       string    `json:"vocab" validate:"required"`
	Kana        string    `json:"kana" validate:"required"`
	Meaning     string    `json:"meaning" validate:"required"`
	English     string    `json:"english" validate:"required"`
	Sentence    string    `json:"sentence" validate:"required"`
	IsJukujikun *bool     `json:"is_jukujikun,omitempty"`
}

type UpdateCardRequest struct {
	Vocab       string `json:"vocab" validate:"required"`
	Kana        string `json:"kana" validate:"required"`
	Meaning     string `json:"meaning" validate:"required"`
	English     string `json:"english" validate:"required"`
	Sentence    string `json:"sentence" validate:"required"`
	IsJukujikun *bool  `json:"is_jukujikun,omitempty"`
}

// Response

type CardResponse struct {
	ID          uuid.UUID `json:"id"`
	Vocab       string    `json:"vocab"`
	Kana        string    `json:"kana"`
	English     string    `json:"english"`
	Meaning     string    `json:"meaning"`
	Sentence    string    `json:"sentence"`
	IsJukujikun bool      `json:"is_jukujikun"`
}
