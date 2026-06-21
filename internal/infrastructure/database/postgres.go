package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valenoirs/flashcard-api/internal/config"
)

type PostgresDatabase struct {
	db        *pgxpool.Pool
	logger    *slog.Logger
	threshold time.Duration
}

type PostgresQueryExecutor interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

func NewPostgresDatabase(cfg *config.Config, log *slog.Logger) (*PostgresDatabase, func(), error) {
	log.Info("establishing postgres database connection pool...")
	dsn := fmt.Sprintf("%s&search_path=%s", cfg.Postgres.ConnectionString, cfg.Postgres.SchemaName)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse database database: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.Postgres.MaxOpenConnection)
	poolConfig.MinConns = int32(cfg.Postgres.MaxIdleConnection)
	poolConfig.MaxConnIdleTime = cfg.Postgres.MaxIdleTime
	poolConfig.MaxConnLifetime = cfg.Postgres.MaxLifetime

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed establishing database connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	cleanup := func() {
		log.Info("closing database connection pool...")
		pool.Close()
	}

	return &PostgresDatabase{
		db:        pool,
		logger:    log,
		threshold: cfg.Postgres.SlowQueryThreshold,
	}, cleanup, nil
}

func (s *PostgresDatabase) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	start := time.Now()
	result, err := s.db.Exec(ctx, query, args...)
	s.logSlowQuery(ctx, start, query, args)
	return result, err
}

func (s *PostgresDatabase) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	start := time.Now()
	rows, err := s.db.Query(ctx, query, args...)
	s.logSlowQuery(ctx, start, query, args)
	return rows, err
}

func (s *PostgresDatabase) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	start := time.Now()
	row := s.db.QueryRow(ctx, query, args...)
	s.logSlowQuery(ctx, start, query, args)
	return row
}

func (s *PostgresDatabase) logSlowQuery(ctx context.Context, start time.Time, query string, args []any) {
	duration := time.Since(start)
	if duration > s.threshold {
		s.logger.WarnContext(ctx, "slow database query row detected",
			slog.String("query", query),
			slog.Any("args", args),
			slog.Duration("duration", duration),
		)
	}
}
