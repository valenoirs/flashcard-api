package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/valenoirs/flashcard-api/internal/domain"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/database"
)

type Card struct {
	ID        uuid.UUID `db:"id"`
	DeckID    uuid.UUID `db:"deck_id"`
	Vocab     string    `db:"vocab"`
	Kana      string    `db:"kana"`
	Meaning   string    `db:"meaning"`
	English   string    `db:"english"`
	Sentence  string    `db:"sentence"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c *Card) ToDomain() *domain.Card {
	return &domain.Card{
		ID:        c.ID,
		DeckID:    c.DeckID,
		Vocab:     c.Vocab,
		Kana:      c.Kana,
		Meaning:   c.Meaning,
		English:   c.English,
		Sentence:  c.Sentence,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func NewCardModel(c *domain.Card) *Card {
	return &Card{
		ID:        c.ID,
		DeckID:    c.DeckID,
		Vocab:     c.Vocab,
		Kana:      c.Kana,
		Meaning:   c.Meaning,
		English:   c.English,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

type cardPostgresAdapter struct {
	db database.PostgresQueryExecutor
}

// GetCardByID implements [domain.CardRepository].
func (c *cardPostgresAdapter) GetCardByID(ctx context.Context, id uuid.UUID) (*domain.Card, error) {
	query := `
	SELECT
		id,
		deck_id,
		vocab,
		kana,
		meaning,
		english,
		sentence,
		created_at,
		updated_at
	FROM cards
	WHERE id = $1
	`

	rows, err := c.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	// defer removed as pgx auto close db connection

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Card])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return result.ToDomain(), nil
}

// CreateCard implements [domain.CardRepository].
func (c *cardPostgresAdapter) CreateCard(ctx context.Context, card *domain.Card) error {
	m := NewCardModel(card)
	query := `
	INSERT INTO cards (id, deck_id, vocab, kana, sentence, meaning, english)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := c.db.Exec(ctx, query,
		m.ID,
		m.DeckID,
		m.Vocab,
		m.Kana,
		m.Sentence,
		m.Meaning,
		m.English,
	)

	return err
}

// DeleteCard implements [domain.CardRepository].
func (c *cardPostgresAdapter) DeleteCard(ctx context.Context, cardID uuid.UUID) error {
	query := "DELETE FROM cards WHERE id = $1"
	commandTag, err := c.db.Exec(ctx, query, cardID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// GetCardList implements [domain.CardRepository].
func (c *cardPostgresAdapter) GetCardList(ctx context.Context, deckID uuid.UUID) ([]domain.Card, error) {
	query := `
	SELECT
		id,
		deck_id,
		vocab,
		kana,
		meaning,
		english,
		sentence,
		created_at,
		updated_at
	FROM cards
	WHERE deck_id = $1
	ORDER BY created_at
	`

	rows, err := c.db.Query(ctx, query, deckID)
	if err != nil {
		return nil, err
	}
	// defer removed as pgx auto close db connection

	results, err := pgx.CollectRows(rows, pgx.RowToStructByName[Card])
	if err != nil {
		return nil, err
	}

	result := make([]domain.Card, len(results))
	for i := range results {
		m := results[i]
		result[i] = *m.ToDomain()
	}

	return result, nil
}

// UpdateCard implements [domain.CardRepository].
func (c *cardPostgresAdapter) UpdateCard(ctx context.Context, card *domain.Card) error {
	m := NewCardModel(card)
	query := `
	UPDATE cards
	SET vocab = $1, kana = $2, meaning = $3, english = $4, sentence = $5
	WHERE id = $6
	`
	commandTag, err := c.db.Exec(ctx, query, m.Vocab, m.Kana, m.Meaning, m.English, m.Sentence, m.ID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func NewCardPostgresAdapter(db database.PostgresQueryExecutor) domain.CardRepository {
	return &cardPostgresAdapter{
		db: db,
	}
}
