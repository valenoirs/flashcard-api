package card

import "github.com/google/uuid"

type CardSentenceCommand struct {
	Position int
	Text     string
	Reading  string
	IsTarget bool
}

type CreateCardCommand struct {
	ID        uuid.UUID
	DeckID    uuid.UUID
	Vocab     string
	Kana      string
	Meaning   string
	English   string
	Sentences []*CardSentenceCommand
}

type UpdateCardCommand struct {
	ID        uuid.UUID
	Vocab     string
	Kana      string
	Meaning   string
	English   string
	Sentences []*CardSentenceCommand
}

type DeleteCardCommand struct {
	ID uuid.UUID
}
