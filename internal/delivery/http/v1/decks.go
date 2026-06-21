package v1

import (
	"encoding/json"
	"net/http"

	"github.com/valenoirs/flashcard-api/internal/delivery/http/v1/request"
	"github.com/valenoirs/flashcard-api/internal/delivery/http/v1/response"
	"github.com/valenoirs/flashcard-api/internal/domain"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/cqrs/command"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/cqrs/query"
	"github.com/valenoirs/flashcard-api/internal/usecase/deck"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/validator"
)

type DeckHandler struct {
	commandRegistry *command.Registry
	queryRegistry   *query.Registry
}

func RegisterDeckRoutes(
	mux *http.ServeMux,
	commandRegistry *command.Registry,
	queryRegistry *query.Registry,
) {
	h := &DeckHandler{
		commandRegistry: commandRegistry,
		queryRegistry:   queryRegistry,
	}

	RouterGroup(mux, "/api/v1/decks", []func(http.Handler) http.Handler{}, func(group func(string, string, http.HandlerFunc)) {
		group(http.MethodPost, "/", h.CreateDeck)
		group(http.MethodGet, "/", h.GetDeckList)
	})
}

func (d *DeckHandler) CreateDeck(w http.ResponseWriter, r *http.Request) {
	req, err := validator.ValidateBody[request.CreateDeckRequest](r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{"errors": err})
		return
	}

	if err := command.Dispatch(r.Context(), d.commandRegistry, req.ToCommand()); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (d *DeckHandler) GetDeckList(w http.ResponseWriter, r *http.Request) {
	data, err := query.Dispatch[*deck.GetDeckListQuery, []domain.Deck](r.Context(), d.queryRegistry, &deck.GetDeckListQuery{})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": "internal server error"})
		return
	}

	res := make([]response.GetDeckResponse, 0, len(data))

	for _, item := range data {
		res = append(res, response.NewGetDeckResponse(&item))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"data": res})
}
