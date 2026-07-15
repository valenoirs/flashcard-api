package card

import "github.com/google/uuid"

type CreateCardCommand struct {
	ID       uuid.UUID
	DeckID   uuid.UUID
	Vocab    string
	Kana     string
	Meaning  string
	English  string
	Sentence string
	IsJukujikun bool
}

type UpdateCardCommand struct {
	ID       uuid.UUID
	Vocab    string
	Kana     string
	Meaning  string
	English  string
	Sentence string
	IsJukujikun bool
}

type DeleteCardCommand struct {
	ID uuid.UUID
}
