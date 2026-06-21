package query

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
)

type Registry struct {
	handlers map[reflect.Type]any
	logger *slog.Logger
}

func NewRegistry(logger *slog.Logger) *Registry{
	return &Registry{
		handlers: make(map[reflect.Type]any),
		logger: logger,
	}
}

func Register[Q Query, R Response](r *Registry, h Handler[Q, R]) {
	var q Q
	qType := reflect.TypeOf(q)
	r.handlers[qType] = h
	r.logger.Info("query handler registered", slog.String("query", qType.String()))
}

func Dispatch[Q Query, R Response](ctx context.Context, r *Registry, q Q) (R, error){
	var zero R
	qType := reflect.TypeOf(q)

	hAny, exists := r.handlers[qType]
	if !exists {
		return zero, fmt.Errorf("no handler registered for command: %s", qType.String())
	}

	h, ok := hAny.(Handler[Q, R])
	if !ok {
		return zero, fmt.Errorf("handler type mismatch for command %s", qType.String())
	}

	return h.Handle(ctx, q)
}
