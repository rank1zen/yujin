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
	}, nil
}

type SummonerProfile struct{}

type ProfilePage struct {
	SummonerName  string
	SummonerLevel int
	Matchlist     []*MatchPlayer
	ProfileIconId int
}

func (p ProfilePage) ProfilePageUrl() string {
	return "https://static.bigbrain.gg/assets/lol/riot_static/14.10.1/img/champion/Jhin.png"
}

func (p ProfilePage) IsRanked() bool {
	return false
}

func (db *DB) ProfilePage(ctx context.Context, puuid string) (*ProfilePage, error) {
	page := new(ProfilePage)

	err := db.pool.QueryRow(ctx, `
	SELECT
		profile_icon_id, summoner_level, 'Temp Name'
	FROM summoner_records_newest
	WHERE puuid = $1;
	`, puuid).Scan(page.ProfileIconId, page.SummonerLevel, page.SummonerName)
	if err != nil {
		return nil, fmt.Errorf("no %w", err)
	}

	matchlist, err := db.UpdateMatchHistory(ctx, puuid)
	if err != nil {
		return nil, err
	}

	page.Matchlist = matchlist

	return page, nil
}

func (db *DB) Health(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *DB) Close() {
	db.pool.Close()
}
