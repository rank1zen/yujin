package postgresql

import (
	"context"

	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/postgresql/db"
)

type SummonerDA struct {
	q *db.Queries
}

func NewSummonerDA(d db.DBTX) *SummonerDA {
	return &SummonerDA{q: db.New(d)}
}

func (s *SummonerDA) Create(ctx context.Context, summoner internal.Summoner) (uint, error) {
	return nil, nil
}

func (s *SummonerDA) Find(ctx context.Context, puuid string) ([]internal.Summoner, error) {
	return nil, nil
}

func (s *SummonerDA) Newest(ctx context.Context, puuid string) (internal.Summoner, error) {
	return nil, nil
}

func (s * SummonerDA) FindByName(ctx context.Context, name string) ([]internal.Summoner, error) {
	return nil, nil
}

func (s * SummonerDA) Delete(ctx context.Context, puuid string) error {
	return nil
}

