package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Queries struct {
	db *pgxpool.Pool
}

func NewQueries(pool *pgxpool.Pool) *Queries {
	return &Queries{db: pool}
}

type SummonerProfileArg struct {
	Name       string
	Puuid      string
	AccountId  string
	SummonerId string
}

func (q *Queries) UpsertSummonerProfile(ctx context.Context, r *SummonerProfileArg) (error) {
	query := `
	INSERT INTO summoner_profile (name, puuid, account_id, summoner_id)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (puuid)
	DO UPDATE SET name = $1, account_id = $2, summoner_id = $4
	`

	_, err := q.db.Exec(ctx, query, r.Name, r.Puuid, r.AccountId, r.SummonerId)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}

	return nil
}
