package service

import (
	"context"

	"github.com/rank1zen/yujin/internal"
)



type SummonerRepo interface {
	Create(ctx context.Context, sum internal.Summoner) (internal.Summoner, error)
	FindRecent(ctx context.Context, id string) (internal.Summoner, error)
}

type SummonerSearchRepo interface {
	// TODO: How does searching work
	Search(ctx context.Context)
}

type Summoner struct {
	repo SummonerRepo
	search SummonerSearchRepo
}

func NewSummoner(repo SummonerRepo, search SummonerSearchRepo) *Summoner {
	return &Summoner{
		repo: repo,
		search: search,
	}
}

func (s *Summoner) Create(ctx context.Context, summoner internal.Summoner) (internal.Summoner, error) {
	res, err := s.repo.Create(ctx, summoner)

	if err != nil {
		return internal.Summoner{}, nil
	}

	return res, nil
}

func (s *Summoner) Find(ctx context.Context) {}
