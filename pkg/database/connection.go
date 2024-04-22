package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Wrapper for underlying Pgx Pool
type DB struct {
	Pgx *pgxpool.Pool
}

// Close DB connection
func (c *DB) Close(ctx context.Context) {
	c.Pgx.Close()
}

func NewFromEnv(ctx context.Context, cfg *Config) (*DB, error) {
	pgxCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pgxCfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		return conn.Ping(ctx) == nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

        err = pool.Ping(ctx)
        if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
        }

	return &DB{Pgx: pool}, nil
}
