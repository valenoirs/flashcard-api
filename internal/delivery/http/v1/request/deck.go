package request

import "github.com/valenoirs/flashcard-api/internal/usecase/deck"

type CreateDeckRequest struct {
	Name string `json:"name" validate:"required"`
}

func (c *CreateDeckRequest) ToCommand() *deck.CreateDeckCommand{
	return &deck.CreateDeckCommand{
		Name: c.Name,
	}
}
