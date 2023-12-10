package riot

import (
	"context"

	"github.com/KnutZuidema/golio/riot/lol"
)

type SummonerQ struct {
	q *lol.SummonerClient
}

func NewSummonerQ(q *lol.SummonerClient) *SummonerQ {
	return &SummonerQ{
		q: q,
	}
}

func (s *SummonerQ) ByPuuid(ctx context.Context, puuid string) (*lol.Summoner, error) {
	r, err := s.q.GetByPUUID(puuid)
	return r, err
}

func (s *SummonerQ) ByName(ctx context.Context, name string) (*lol.Summoner, error) {
	r, err := s.q.GetByName(name)
	return r, err
}
