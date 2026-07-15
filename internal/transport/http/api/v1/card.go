package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/valenoirs/flashcard-api/internal/domain"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/cqrs/command"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/cqrs/query"
	"github.com/valenoirs/flashcard-api/internal/transport/http/api/v1/dto"
	"github.com/valenoirs/flashcard-api/internal/transport/http/router"
	"github.com/valenoirs/flashcard-api/internal/transport/http/validator"
	"github.com/valenoirs/flashcard-api/internal/usecase/card"
)

type CardHandler struct {
	commandRegistry *command.Registry
	queryRegistry   *query.Registry
}

func RegisterCardRoutes(mux *http.ServeMux, commandRegistry *command.Registry, queryRegistry *query.Registry) {
	h := &CardHandler{
		commandRegistry: commandRegistry,
		queryRegistry:   queryRegistry,
	}

	router.Group(mux, "/api/v1/cards", []func(http.Handler) http.Handler{}, func(group func(string, string, http.HandlerFunc)) {
		group(http.MethodPost, "/", h.CreateCard)
		group(http.MethodPut, "/{id}", h.UpdateCard)
		group(http.MethodDelete, "/{id}", h.DeleteCard)
		group(http.MethodGet, "/{deckID}", h.GetCardList)
	})
}

func (c *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	req, valErr := validator.ValidateBody[dto.CreateCardRequest](r.Body)
	if valErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{"errors": valErr})
		return
	}

	cardID := uuid.Must(uuid.NewV7())

	cmd := &card.CreateCardCommand{
		ID:       cardID,
		DeckID:   req.DeckID,
		Vocab:    req.Vocab,
		Kana:     req.Kana,
		Meaning:  req.Meaning,
		English:  req.English,
		Sentence: req.Sentence,
	}

	if err := command.Dispatch(r.Context(), c.commandRegistry, cmd); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	q := &card.GetCardDetailQuery{
		ID: cardID,
	}

	data, err := query.Dispatch[*card.GetCardDetailQuery, *domain.Card](r.Context(), c.queryRegistry, q)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	// NOTE: might cause not found error when switch to read/write db
	res := toCardResponse(data)

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

	req, err := validator.ValidateBody[dto.UpdateCardRequest](r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{"errors": err})
		return
	}

	cmd := &card.UpdateCardCommand{
		ID:       cardID,
		Vocab:    req.Vocab,
		Kana:     req.Kana,
		Meaning:  req.Meaning,
		English:  req.English,
		Sentence: req.Sentence,
	}

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

	cmd := &card.DeleteCardCommand{
		ID: cardID,
	}

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

	q := &card.GetCardListQuery{
		DeckID: deckID,
	}

	data, err := query.Dispatch[*card.GetCardListQuery, []domain.Card](r.Context(), c.queryRegistry, q)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	res := make([]dto.CardResponse, len(data))

	for i := range data {
		d := data[i]
		res[i] = toCardResponse(&d)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"data": res})
}

func toCardResponse(c *domain.Card) dto.CardResponse {
	parsedSentence := parseJapaneseMarkdown(c.Sentence)
	return dto.CardResponse{
		ID:        c.ID,
		Vocab:     c.Vocab,
		Kana:      c.Kana,
		Meaning:   c.Meaning,
		English:   c.English,
		Sentence:  parsedSentence,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func parseJapaneseMarkdown(parsedText string) string {
	htmlResult := parsedText

	// 1. Handle Furigana: [難:むずか] -> <ruby>難<rt>むずか</rt></ruby>
	furiganaRegex := regexp.MustCompile(`\[([^:]+):([^\]]+)\]`)
	htmlResult = furiganaRegex.ReplaceAllString(htmlResult, "<ruby>$1<rt>$2</rt></ruby>")

	// 2. Handle Target/Bold Words: **一番** -> <strong class="target-word">一番</strong>
	boldRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	htmlResult = boldRegex.ReplaceAllString(htmlResult, `<strong>$1</strong>`)

	return htmlResult
}
