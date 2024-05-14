package database

import (
	"context"
	"fmt"

	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/pkg/logging"
)

var (
	soloQueueType = 420
	soloqOption   = lol.MatchListOptions{Queue: &soloQueueType}
)

// This is a wrapper for exclusivly pgx "QUERY" logic
type pgxDB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...any) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Close()
}

// DB represents a GIGA interface for accessing repository
type DB interface {
	Summoner() SummonerQuery
	League() LeagueQuery
	Match() MatchQuery

	Close()
}

type db struct {
	summoner SummonerQuery
	league   LeagueQuery
	match    MatchQuery

	pgx pgxDB
}

// NewDB creates and returns a new database from string
func NewDB(ctx context.Context, url string) (DB, error) {
	log := logging.FromContext(ctx).Sugar()

	pgxCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres connection string: %w", err)
	}

	pgxCfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		return conn.Ping(ctx) == nil
	}

	log.Infof("attempting to connect to postgres...")
	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return &db{
		summoner: NewSummonerQuery(pool),
		league:   NewLeagueQuery(pool),
		match:    NewMatchQuery(pool),
		pgx:      pool,
	}, nil
}

// WithNewDB creates a new database from string and attaches it to some allowing interface
func WithNewDB(ctx context.Context, e interface{ SetDatabase(DB) }, url string) error {
	db, err := NewDB(ctx, url)
	if err != nil {
		return err
	}

	e.SetDatabase(db)
	return nil
}

func (d *db) Summoner() SummonerQuery { return d.summoner }
func (d *db) League() LeagueQuery     { return d.league }
func (d *db) Match() MatchQuery       { return d.match }

// Close closes the DB connection
func (d *db) Close() {
	d.pgx.Close()
}
