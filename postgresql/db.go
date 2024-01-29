package postgresql

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"go.uber.org/zap"
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

	err = migrator.LoadMigrations(os.DirFS("../db/migrations"))
	if err != nil {
		return fmt.Errorf("could not load migrations: %w", err)
	}

	err = migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("could not migrate: %w", err)
	}

	return nil
}

func BackoffRetryPool(ctx context.Context, url string, log *zap.Logger) (*pgxpool.Pool, error) {
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

	notify := func(err error, d time.Duration) {
		log.Warn("could not connect to postgres", zap.Error(err))
	}

	if err := backoff.RetryNotify(op, backoff.WithMaxRetries(b, 5), notify); err != nil {
		return nil, err
	}

	return pool, nil
}

func CheckPool(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return errors.New("no database connection")
	}
	return pool.Ping(ctx)
}
