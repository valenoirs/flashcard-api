package command

import "context"

type Command any

type Handler[C Command] interface {
	Handle(ctx context.Context, cmd C) error
}
