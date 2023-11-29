package service

import (
	"context"

	"github.com/rank1zen/yujin/internal"
)

type SummonerRepo interface {
	Create(ctx context.Context, params internal.PSummonerCreate) (internal.Summoner, error)
	Find(ctx context.Context, params internal.PSummonerFind) ([]internal.Summoner, error)
	Newest(ctx context.Context, params internal.PSummonerNewest) (internal.Summoner, error)
}

type SummonerSearchRepo interface {
	// TODO: How does searching work
	Search(ctx context.Context)
}

type Summoner struct {
	repo   SummonerRepo
	search SummonerSearchRepo
}

func NewSummoner(repo SummonerRepo, search SummonerSearchRepo) *Summoner {
	return &Summoner{
		repo:   repo,
		search: search,
	}
}

func (s *Summoner) Create(ctx context.Context, params internal.PSummonerCreate) (internal.Summoner, error) {
	res, err := s.repo.Create(ctx, params)

	if err != nil {
		return internal.Summoner{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo.Create")
	}

	return res, nil
}

func (s *Summoner) Find(ctx context.Context, params internal.PSummonerFind) ([]internal.Summoner, error) {
	res, err := s.repo.Find(ctx, params)

	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo.Find")
	}

	return res, nil
}

func (s *Summoner) Newest(ctx context.Context, params internal.PSummonerNewest) (internal.Summoner, error) {
	res, err := s.repo.Newest(ctx, params)

	if err != nil {
		return internal.Summoner{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo.Newest")
	}

	return res, nil
}

