package service

import (
	"context"

	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/postgresql"
	"github.com/rank1zen/yujin/internal/riot"
)

type SummonerService struct {
	repo           *postgresql.SummonerDA
	summonerClient riot.SummonerQueryClient
}

func NewSummonerService(repo *postgresql.SummonerDA) *SummonerService {
	return &SummonerService{
		repo: repo,
	}
}

func (s *SummonerService) Create(ctx context.Context, summoner internal.SummonerWithIds) (string, error) {
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

func (s *SummonerService) QueryRiot(ctx context.Context, name string) (internal.SummonerWithIds, error) {
	r, err := s.summonerClient.ByName(ctx, &riot.SummonerByNameRequest{})
	if err != nil {
		return internal.SummonerWithIds{}, err
	}
	return internal.Summoner{

		Name: r.GetName(),
	}, nil
}
