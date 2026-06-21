package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/valenoirs/flashcard-api/internal/delivery/http/v1/request"
	"github.com/valenoirs/flashcard-api/internal/delivery/http/v1/response"
	"github.com/valenoirs/flashcard-api/internal/domain"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/cqrs/command"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/cqrs/query"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/validator"
	"github.com/valenoirs/flashcard-api/internal/usecase/card"
)

type CardHandler struct {
	commandRegistry *command.Registry
	queryRegistry   *query.Registry
}

func RegisterCardRoutes(
	mux *http.ServeMux,
	commandRegistry *command.Registry,
	queryRegistry *query.Registry,
) {
	h := &CardHandler{
		commandRegistry: commandRegistry,
		queryRegistry:   queryRegistry,
	}

	RouterGroup(mux, "/api/v1/cards", []func(http.Handler) http.Handler{}, func(group func(string, string, http.HandlerFunc)) {
		group(http.MethodPost, "/", h.CreateCard)
		group(http.MethodPut, "/{id}", h.UpdateCard)
		group(http.MethodDelete, "/{id}", h.DeleteCard)
		group(http.MethodGet, "/{deckID}", h.GetCardList)
	})
}

func (c *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	req, valErr := validator.ValidateBody[request.CreateCardRequest](r.Body)
	if valErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{"errors": valErr})
		return
	}

	cardID, err := uuid.NewV7()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": "failed to generate uuid"})
		return
	}

	cmd := request.NewCreateCardCommand(req, cardID)

	if err = command.Dispatch(r.Context(), c.commandRegistry, cmd); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	q := request.NewGetCardDetailQuery(cardID)

	data, err := query.Dispatch[*card.GetCardDetailQuery, *domain.Card](r.Context(), c.queryRegistry, q)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	res := response.NewGetCardResponse(data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"data": res})
}

func (c *CardHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	cardIDstr := r.PathValue("id")

	cardID, parseErr := uuid.Parse(cardIDstr)
	if parseErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{"errors": "invalid card id"})
		return
	}

	req, err := validator.ValidateBody[request.UpdateCardRequest](r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{"errors": err})
		return
	}

	cmd := request.NewUpdateCardCommand(req, cardID)

	if err := command.Dispatch(r.Context(), c.commandRegistry, cmd); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *CardHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	cardIDstr := r.PathValue("id")

	cardID, err := uuid.Parse(cardIDstr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{"errors": err})
		return
	}

	cmd := request.NewDeleteCardCommand(cardID)

	if err := command.Dispatch(r.Context(), c.commandRegistry, cmd); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]any{"error": "card not found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *CardHandler) GetCardList(w http.ResponseWriter, r *http.Request) {
	deckIDstr := r.PathValue("deckID")

	deckID, err := uuid.Parse(deckIDstr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{"errors": err.Error()})
		return
	}

	q := request.NewGetCardListQuery(deckID)

	data, err := query.Dispatch[*card.GetCardListQuery, []domain.Card](r.Context(), c.queryRegistry, q)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	res := make([]response.GetCardResponse, 0, len(data))

	for _, item := range data {
		res = append(res, response.NewGetCardResponse(&item))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"data": res})
}
