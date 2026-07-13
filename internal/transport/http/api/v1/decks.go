package v1

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/valenoirs/flashcard-api/internal/domain"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/cqrs/command"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/cqrs/query"
	"github.com/valenoirs/flashcard-api/internal/transport/http/api/v1/dto"
	"github.com/valenoirs/flashcard-api/internal/transport/http/router"
	"github.com/valenoirs/flashcard-api/internal/transport/http/validator"
	"github.com/valenoirs/flashcard-api/internal/usecase/deck"
)

type DeckHandler struct {
	commandRegistry *command.Registry
	queryRegistry   *query.Registry
}

func RegisterDeckRoutes(mux *http.ServeMux, commandRegistry *command.Registry, queryRegistry *query.Registry) {
	h := &DeckHandler{
		commandRegistry: commandRegistry,
		queryRegistry:   queryRegistry,
	}

	router.Group(mux, "/api/v1/decks", []func(http.Handler) http.Handler{}, func(group func(string, string, http.HandlerFunc)) {
		group(http.MethodPost, "/", h.CreateDeck)
		group(http.MethodGet, "/", h.GetDeckList)
	})
}

func (d *DeckHandler) CreateDeck(w http.ResponseWriter, r *http.Request) {
	req, valErr := validator.ValidateBody[dto.CreateDeckRequest](r.Body)
	if valErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{"errors": valErr})
		return
	}

	deckID, err := uuid.NewV7()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": "failed to generate uuid"})
		return
	}

	cmd := &deck.CreateDeckCommand{
		DeckID: deckID,
		Name:   req.Name,
	}

	if err := command.Dispatch(r.Context(), d.commandRegistry, cmd); err != nil {
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

	res := make([]dto.DeckResponse, 0, len(data))

	for _, item := range data {
		res = append(res, toDeckResponse(&item))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"data": res})
}

func toDeckResponse(d *domain.Deck) dto.DeckResponse {
	return dto.DeckResponse{
		ID:        d.ID,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
