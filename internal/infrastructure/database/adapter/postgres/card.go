package postgres

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/valenoirs/flashcard-api/internal/domain"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/database"
)

type Card struct {
	ID        uuid.UUID
	DeckID    uuid.UUID
	Front     string
	Back      string
	Note      *string
	Class     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Card) ToDomain() *domain.Card {
	return &domain.Card{
		ID:        c.ID,
		DeckID:    c.DeckID,
		Front:     c.Front,
		Back:      c.Back,
		Note:      c.Note,
		Class:     c.Class,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func NewCardModel(c *domain.Card) *Card {
	return &Card{
		ID:        c.ID,
		DeckID:    c.DeckID,
		Front:     strings.ToLower(c.Front),
		Back:      strings.ToLower(c.Back),
		Note:      c.Note,
		Class:     c.Class,
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
	SELECT id, front, back, note, class, created_at, updated_at
	FROM cards
	WHERE id = $1
	`

	var result Card

	err := c.db.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.Front,
		&result.Back,
		&result.Note,
		&result.Class,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
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
	INSERT INTO cards (id, deck_id, front, back, note, class)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := c.db.Exec(ctx, query,
		m.ID,
		m.DeckID,
		m.Front,
		m.Back,
		m.Note,
		m.Class,
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
	SELECT id, front, back, note, class, created_at, updated_at
	FROM cards
	WHERE deck_id = $1
	ORDER BY created_at
	`

	rows, err := c.db.Query(ctx, query, deckID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Card

	for rows.Next() {
		var m Card

		err := rows.Scan(&m.ID, &m.Front, &m.Back, &m.Note, &m.Class, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, *m.ToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateCard implements [domain.CardRepository].
func (c *cardPostgresAdapter) UpdateCard(ctx context.Context, card *domain.Card) error {
	m := NewCardModel(card)
	query := `
	UPDATE cards
	SET front = $1, back = $2, note = $3, class = $4
	WHERE id = $5
	`
	commandTag, err := c.db.Exec(ctx, query, m.Front, m.Back, m.Note, m.Class, m.ID)
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
