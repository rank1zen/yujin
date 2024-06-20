package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/pkg/logging"
)

// This is a wrapper for exclusivly pgx "QUERY" logic
type pgxDB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...any) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

type zapTracer struct{}

func (z *zapTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	logger := logging.FromContext(ctx).Sugar()
	logger.Debugf("Executing SQL: %s", data.SQL)
	return ctx
}

func (z *zapTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	logger := logging.FromContext(ctx).Sugar()
	logger.Debugf("The flip is a command tag: %v", data.CommandTag)
}

type DB struct {
	pgx *pgxpool.Pool
}

// NewDB creates and returns a new database from string
func NewDB(ctx context.Context, url string) (*DB, error) {
	pgxCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres connection string: %w", err)
	}

	pgxCfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		return conn.Ping(ctx) == nil
	}

	tracer := zapTracer{}
	pgxCfg.ConnConfig.Tracer = &tracer

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return &DB{pgx: pool}, nil
}

// WithNewDB creates a new database from string and attaches it to some allowing interface
func WithNewDB(ctx context.Context, e interface{ SetDatabase(*DB) }, url string) error {
	db, err := NewDB(ctx, url)
	if err != nil {
		return err
	}

	e.SetDatabase(db)
	return nil
}

func (d *DB) Health(ctx context.Context) error {
	return d.pgx.Ping(ctx)
}

func (d *DB) Close() {
	d.pgx.Close()
}
