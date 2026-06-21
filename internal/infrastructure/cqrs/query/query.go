package query

import "context"

type Query any

type Response any

type Handler[Q Query, R Response] interface {
	Handle(ctx context.Context, query Q) (R, error)
}
