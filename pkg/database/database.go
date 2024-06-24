package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/pkg/riot"
	"go.opentelemetry.io/otel/trace"
)

// DB represents all services for this web app
type DB struct {
	pool    *pgxpool.Pool
	riot    *riot.Client
	tracer  trace.Tracer
	service *service
}

// NewDB creates and returns a new database from string
func NewDB(ctx context.Context, url string) (*DB, error) {
	pool, err := newPgxPool(ctx, url)
	if err != nil {
		return nil, err
	}

	riot := riot.NewClient()

	return &DB{
		pool: pool,
		riot: riot,
		service: &service{
			riot: riot,
		},
	}, nil
}

type GetMatch func(ctx context.Context, puuid string, page int) ([]MatchPlayer, error)

// GetMatchHistory gets from DB, the first 5 recent matches available
func (db *DB) GetMatchHistory(ctx context.Context, puuid string, page int) ([]MatchPlayer, error) {
	pagesize := 5
	return db.service.getPlayerMatchHstory(ctx, db.pool, puuid, 5*page, pagesize)
}

// UpdateMatchHistory fetches and inserts the first 20 matches
func (db *DB) UpdateMatchHistory(ctx context.Context, puuid string) error {
	ids, err := db.service.fetchNewMatches(ctx, db.pool, puuid, 0, 20)
	if err != nil {
		return fmt.Errorf("fetch matches: %w", err)
	}

	_, err = db.service.insertMatches(ctx, db.pool, ids)
	if err != nil {
		return err
	}

	return nil
}

// TODO: implement
func (db *DB) GetSummonerProfile(ctx context.Context, puuid string) (*SummonerProfile, error) {
	return nil, nil
}

// TODO: implement
func (db *DB) FetchEntireMatchHistory(ctx context.Context, puuid string) error {
	return nil
}

func (db *DB) Health(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *DB) Close() {
	db.pool.Close()
}
