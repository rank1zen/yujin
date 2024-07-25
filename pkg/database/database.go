package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/pkg/riot"
	"go.opentelemetry.io/otel/trace"
)

type DB struct {
	pool   *pgxpool.Pool
	riot   *riot.Client
	tracer trace.Tracer
}

func NewDB(ctx context.Context, url string) (*DB, error) {
	pgxCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres connection string: %w", err)
	}

	pgxCfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		return conn.Ping(ctx) == nil
	}

	pgxCfg.ConnConfig.Tracer = &tracer{}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	riot := riot.NewClient()

	return &DB{
		pool: pool,
		riot: riot,
	}, nil
}

func (db *DB) Health(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *DB) Close() {
	db.pool.Close()
}
