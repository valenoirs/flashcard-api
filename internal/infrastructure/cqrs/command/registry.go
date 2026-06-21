package command

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
)

type Registry struct {
	handlers map[reflect.Type]any
	logger   *slog.Logger
}

func NewRegistry(logger *slog.Logger) *Registry {
	return &Registry{
		handlers: make(map[reflect.Type]any),
		logger:   logger,
	}
}

func Register[C Command](r *Registry, h Handler[C]) {
	var c C
	cType := reflect.TypeOf(c)
	r.handlers[cType] = h
	r.logger.Info("command handler registered", slog.String("command", cType.String()))
}

func Dispatch[C Command](ctx context.Context, r *Registry, c C) error {
	cType := reflect.TypeOf(c)

	hAny, exists := r.handlers[cType]
	if !exists {
		return fmt.Errorf("no handler registered for command: %s", cType.String())
	}

	h, ok := hAny.(Handler[C])
	if !ok {
		return fmt.Errorf("handler type mismatch for command %s", cType.String())
	}

	return h.Handle(ctx, c)
}
