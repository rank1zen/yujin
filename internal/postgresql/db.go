package postgresql

import (
	"context"
	"fmt"
	"os"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
)

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}

	migrator, err := migrate.NewMigrator(ctx, conn.Conn(), "public.schema_version")
	if err != nil {
		return fmt.Errorf("could not create migrator: %w", err)
	}

	err = migrator.LoadMigrations(os.DirFS("../../migrations"))
	if err != nil {
		return fmt.Errorf("could not load migrations: %w", err)
	}

	err = migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("could not migrate: %w", err)
	}

	return nil
}

func NewBackoffPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool

	op := func() error {
		var err error
		pool, err = pgxpool.New(ctx, url)
		if err != nil {
			return err
		}
		return pool.Ping(ctx)
	}

	b := backoff.NewExponentialBackOff()

	if err := backoff.Retry(op, backoff.WithMaxRetries(b, 5)); err != nil {
		return nil, err
	}

	return pool, nil
}

func NewConnectionPool(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(cfg.ConnectionURL())
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	conf.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		return conn.Ping(ctx) == nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return pool, nil
}
