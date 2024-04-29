package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SummonerV4Query interface {
        GetSummonerRecords(context.Context, ...RecordFilter) ([]*SummonerRecord, error)
        CountSummonerRecords(context.Context, ...RecordFilter) (int64, error)
        InsertSummonerRecords(context.Context, []*SummonerRecord) (int64, error)
        DeleteSummonerRecords(context.Context) error
}

type LeagueV4Query interface {
        GetLeagueRecords(context.Context)
        CountLeagueRecords(context.Context)
        InsertLeagueRecords(context.Context)
}

type MatchV5Query interface {
        InsertMatchRecords(context.Context, []*MatchRecord) (int64, error)
        InsertMatchTeamRecords(context.Context, []*MatchTeamRecord) (int64, error)
        InsertMatchBanRecords(context.Context, []*MatchBanRecord) (int64, error)
        InsertMatchParticipantRecords(context.Context, []*MatchParticipantRecord) (int64, error)

        GetMatchRecords(context.Context, ...RecordFilter) ([]*MatchRecord, error)
        GetMatchRecordsMatchingIds(context.Context, []string) ([]*MatchRecord, []string, error)
        GetMatchlist(context.Context, string) ([]string, error)
}

// This is a wrapper for exclusivly "QUERY" logic
type pgDB interface {
        Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
        CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
        Close()
}

// DB represents a collection of things to access database things
type DB struct {
	SummonerV4 SummonerV4Query
	LeagueV4   LeagueV4Query
	MatchV5    MatchV5Query

	pgx pgDB
	Cfg *Config
}

type Config struct {
        Url string
}

func (c *Config) DSN() string {
        return c.Url
}

func NewConfig(connString string) *Config {
        return &Config{Url: connString}
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

	return &DB{pgx: pool}, nil
}

// Close DB connection
func (c *DB) Close(ctx context.Context) {
	c.pgx.Close()
}

func (c *DB) GetPgx() pgDB {
        return c.pgx
}
