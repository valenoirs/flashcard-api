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
	ID        uuid.UUID  `db:"id"`
	DeckID    uuid.UUID  `db:"deck_id"`
	Vocab     string     `db:"vocab"`
	Kana      string     `db:"kana"`
	Meaning   string     `db:"meaning"`
	English   string     `db:"english"`
	Sentences []Sentence `db:"sentences" json:"sentences"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}

func (c *Card) ToDomain() *domain.Card {
	sentences := make([]*domain.Sentence, len(c.Sentences))

	for i := range c.Sentences {
		sentences[i] = c.Sentences[i].ToDomain()
	}

	return &domain.Card{
		ID:        c.ID,
		DeckID:    c.DeckID,
		Vocab:     c.Vocab,
		Kana:      c.Kana,
		Meaning:   c.Meaning,
		English:   c.English,
		Sentences: sentences,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

type Sentence struct {
	ID       uuid.UUID `db:"id" json:"id"`
	CardID   uuid.UUID `db:"card_id" json:"card_id"`
	Position int       `db:"position" json:"position"`
	Text     string    `db:"text" json:"text"`
	Reading  string    `db:"reading" json:"reading"`
	IsTarget bool      `db:"is_target" json:"is_target"`
}

func (s *Sentence) ToDomain() *domain.Sentence {
	return &domain.Sentence{
		ID:       s.ID,
		CardID:   s.CardID,
		Position: s.Position,
		Text:     s.Text,
		Reading:  s.Reading,
		IsTarget: s.IsTarget,
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
		c.id,
		c.deck_id,
		c.vocab,
		c.kana,
		c.meaning,
		c.english,
		c.created_at,
		c.updated_at,
		COALESCE(
			jsonb_agg(
				jsonb_build_object(
					'id', s.id,
					'card_id', s.card_id,
					'position', s.position,
					'text', s.text,
					'reading', s.reading,
					'is_target', s.is_target
				) ORDER BY s.position ASC
			) FILTER (WHERE s.id IS NOT NULL), 
			'[]'::jsonb
		) as sentences
	FROM cards AS c
	LEFT JOIN sentences AS s ON s.card_id = c.id
	WHERE c.id = $1
	GROUP BY c.id
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
	INSERT INTO cards (id, deck_id, vocab, kana, meaning, english)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := c.db.Exec(ctx, query,
		m.ID,
		m.DeckID,
		m.Vocab,
		m.Kana,
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
		c.id,
		c.deck_id,
		c.vocab,
		c.kana,
		c.meaning,
		c.english,
		c.created_at,
		c.updated_at,
		COALESCE(
			jsonb_agg(
				jsonb_build_object(
					'id', s.id,
					'card_id', s.card_id,
					'position', s.position,
					'text', s.text,
					'reading', s.reading,
					'is_target', s.is_target
				) ORDER BY s.position ASC
			) FILTER (WHERE s.id IS NOT NULL), 
			'[]'::jsonb
		) as sentences
	FROM cards AS c
	LEFT JOIN sentences AS s ON s.card_id = c.id
	WHERE c.deck_id = $1
	GROUP BY c.id
	ORDER BY c.created_at;
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
	SET vocab = $1, kana = $2, meaning = $3, english = $4
	WHERE id = $5
	`
	commandTag, err := c.db.Exec(ctx, query, m.Vocab, m.Kana, m.Meaning, m.English, m.ID)
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
