package service

import (
	"context"

	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/postgresql"
)

type SummonerService struct {
	repo *postgresql.SummonerDA
}

func NewSummonerService(repo *postgresql.SummonerDA) *SummonerService {
	return &SummonerService{
		repo: repo,
	}
}

func (s *SummonerService) Create(ctx context.Context, summoner internal.Summoner) (string, error){
	id, err := s.repo.Create(ctx, summoner)
	if err != nil {
		return id, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo.Create")
	}
	return id, nil
}

func (s *SummonerService) FindByPuuid(ctx context.Context, puuid string) ([]internal.Summoner, error) {
	r, err := s.repo.Find(ctx, puuid)
	if err != nil {
		return []internal.Summoner{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo.Find")
	}

	return r, nil
}

func (s *SummonerService) NewestByPuuid(ctx context.Context, puuid string) (internal.Summoner, error) {
	r, err := s.repo.Newest(ctx, puuid)
	if err != nil {
		return internal.Summoner{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo.Newest")
	}

	return r, nil
}
