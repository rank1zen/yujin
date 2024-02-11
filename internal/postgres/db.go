package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"github.com/rank1zen/yujin/internal/logging"
)

func NewBackoffPool(ctx context.Context, connstr string) (*pgxpool.Pool, error) {
	log := logging.FromContext(ctx)

	var pool *pgxpool.Pool

	op := func() error {
		var err error
		pool, err = NewConnectionPool(ctx, connstr)
		if err != nil {
			return err
		}
		return pool.Ping(ctx)
	}

	n := func(err error, d time.Duration) {
		log.Info(err)
	}

	b := backoff.NewExponentialBackOff()

	err := backoff.RetryNotify(op, backoff.WithMaxRetries(b, 5), n)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func NewConnectionPool(ctx context.Context, connstr string) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(connstr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	conf.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		err := RegisterCompositeTypes(ctx, conn)
		if err != nil {
			return false
		}

		return conn.Ping(ctx) == nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return pool, nil
}

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

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

func RegisterCompositeTypes(ctx context.Context, conn *pgx.Conn) error {
	for _, t := range []string{
		"team_champion_ban",
		"team_objective",
	} {
		datatype, err := conn.LoadType(ctx, t)
		if err != nil {
			return err
		}
		conn.TypeMap().RegisterType(datatype)
	}
	return nil
}
