package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/valenoirs/flashcard-api/internal/domain"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/database"
)

type Deck struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (d *Deck) ToDomain() *domain.Deck {
	return &domain.Deck{
		ID:        d.ID,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func NewDeckModel(d *domain.Deck) *Deck {
	return &Deck{
		ID:        d.ID,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

type deckPostgresAdapter struct {
	db database.PostgresQueryExecutor
}

// CreateDeck implements [domain.DeckRepository].
func (d *deckPostgresAdapter) CreateDeck(ctx context.Context, deck *domain.Deck) error {
	m := NewDeckModel(deck)
	query := "INSERT INTO decks (name) VALUES ($1)"
	_, err := d.db.Exec(ctx, query, m.Name)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			fmt.Printf("Postgres Error: %s (Code: %s)\n", pgErr.Message, pgErr.Code)
		}
	}
	return err
}

// DeleteDeck implements [domain.DeckRepository].
func (d *deckPostgresAdapter) DeleteDeck(ctx context.Context, deckID uuid.UUID) error {
	query := "DELETE FROM decks WHERE id = $1"
	_, err := d.db.Exec(ctx, query, deckID)
	return err
}

// GetDeckList implements [domain.DeckRepository].
func (d *deckPostgresAdapter) GetDeckList(ctx context.Context) ([]domain.Deck, error) {
	query := `
	SELECT id, name, created_at, updated_at
	FROM decks
	ORDER BY created_at
	`

	rows, err := d.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domainData []domain.Deck

	for rows.Next() {
		var m Deck
		err := rows.Scan(&m.ID, &m.Name, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		domainData = append(domainData, *m.ToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return domainData, nil
}

// UpdateDeck implements [domain.DeckRepository].
func (d *deckPostgresAdapter) UpdateDeck(ctx context.Context, deck *domain.Deck) error {
	panic("unimplemented")
}

func NewDeckPostgresAdapter(db database.PostgresQueryExecutor) domain.DeckRepository {
	return &deckPostgresAdapter{
		db: db,
	}
}
