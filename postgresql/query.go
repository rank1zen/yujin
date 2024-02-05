package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Query struct {
	LeagueV4   *LeagueV4Query
	SummonerV4 *SummonerV4Query
	db         *pgxpool.Pool
}

func NewQuery(pool *pgxpool.Pool) *Query {
	return &Query{
		LeagueV4: &LeagueV4Query{db: pool},
		SummonerV4: &SummonerV4Query{db: pool},
		db: pool,
	}
}

type SummonerProfileArg struct {
	Name       string
	Puuid      string
	AccountId  string
	SummonerId string
}

func (q *Query) UpsertSummonerProfile(ctx context.Context, r *SummonerProfileArg) error {
	query := `
	INSERT INTO summoner_profile (name, puuid, account_id, summoner_id)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (puuid)
	DO UPDATE SET name = $1, account_id = $2, summoner_id = $4
	`

	_, err := q.db.Exec(ctx, query, r.Name, r.Puuid, r.AccountId, r.SummonerId)
	return err
}
