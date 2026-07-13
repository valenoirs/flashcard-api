package http

import (
	"net/http"

	v1 "github.com/valenoirs/flashcard-api/internal/transport/http/api/v1"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/cqrs/command"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/cqrs/query"
)

func NewHTTPRouter(
	mux *http.ServeMux,
	commandRegistry *command.Registry,
	queryRegistry *query.Registry,
) {
	v1.RegisterCardRoutes(mux, commandRegistry, queryRegistry)
	v1.RegisterDeckRoutes(mux, commandRegistry, queryRegistry)
}
