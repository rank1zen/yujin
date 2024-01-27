package postgresql

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"go.uber.org/zap"
)

func Migrate(ctx context.Context, conn *pgx.Conn, log *zap.Logger) {
	log.Info("migrating database")

	migrator, err := migrate.NewMigrator(ctx, conn, "public.schema_version")
	if err != nil {
		log.Fatal("Could not create migrator: %s", zap.Error(err))
	}

	err = migrator.LoadMigrations(os.DirFS("../db/migrations"))
	if err != nil {
		log.Fatal("Could not load migrations: %s", zap.Error(err))
	}

	err = migrator.Migrate(ctx)
	if err != nil {
		log.Fatal("Could not migrate: %s", zap.Error(err))
	}
}

func BackoffRetryPool(ctx context.Context, url string, log *zap.Logger) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool

	op := func() error {
		pool, err := pgxpool.New(ctx, url)
		if err != nil {
			return err
		}
		return pool.Ping(ctx)
	}

	b := backoff.NewExponentialBackOff()

	notify := func(err error, d time.Duration) {
		log.Warn("could not connect to postgres", zap.Duration("backoff", d))
	}

	log.Info("connecting to postgresql")

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
