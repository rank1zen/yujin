package service

import (
	"context"
	"time"

	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/postgresql"
	"github.com/rank1zen/yujin/internal/riot"
)

type SummonerService struct {
	repo *postgresql.SummonerDA
	riot *riot.SummonerQ
}

func NewSummonerService(repo *postgresql.SummonerDA, riot *riot.SummonerQ) *SummonerService {
	return &SummonerService{
		repo: repo,
		riot: riot,
	}
}

func (s *SummonerService) Create(ctx context.Context, summoner Summoner) (string, error){
	msg, err := s.repo.Create(ctx, summoner.Cast())
	if err != nil {
		return msg, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo.Create")
	}
	return msg, nil
}

func (s *SummonerService) QueryRiot(ctx context.Context, name string) (string, error) {
	sum, err := s.riot.ByName(ctx, name)
	if err != nil {
		return "", internal.WrapErrorf(err, internal.ErrorCodeUnknown, "riot.ByName")
	}
	msg, err := s.Create(ctx, riotCast(sum))
	if err != nil {
		return msg, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "QueryRiot")
	}
	return msg, nil
}

func (s *SummonerService) Search() {}

func (s *SummonerService) NewestSearch() {}

func (s *SummonerService) FindByPuuid(ctx context.Context, puuid string) ([]Summoner, error) {
	r, err := s.repo.Find(ctx, puuid)

	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo.Find")
	}

	res := make([]Summoner, len(r))
	for i, sum := range r {
		res[i] = dbCast(sum)
	}

	return res, nil
}

func (s *SummonerService) NewestByPuuid(ctx context.Context, puuid string) (Summoner, error) {
	r, err := s.repo.Newest(ctx, puuid)

	if err != nil {
		return Summoner{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo.Newest")
	}

	return dbCast(r), nil
}

type Summoner struct {
	Puuid         string
	AccountId     string
	SummonerId    string
	Level         int64
	ProfileIconId int32
	Name          string
	LastRevision  time.Time
	TimeStamp     time.Time
}

func (s Summoner) Cast() postgresql.Summoner {
	return postgresql.Summoner{
		Puuid:         s.Puuid,
		AccountId:     s.AccountId,
		SummonerId:    s.SummonerId,
		Level:         s.Level,
		ProfileIconId: s.ProfileIconId,
		Name:          s.Name,
		LastRevision:  s.LastRevision.Unix(),
		TimeStamp:     s.TimeStamp.Unix(),
	}
}

func dbCast(summoner postgresql.Summoner) Summoner {
	return Summoner{
		Puuid:         summoner.Puuid,
		AccountId:     summoner.AccountId,
		SummonerId:    summoner.SummonerId,
		Level:         summoner.Level,
		ProfileIconId: summoner.ProfileIconId,
		Name:          summoner.Name,
		LastRevision:  time.Unix(summoner.LastRevision, 0),
		TimeStamp:     time.Unix(summoner.TimeStamp, 0),
	}
}

func riotCast(summoner riot.Summoner) Summoner {
	return Summoner{
		Puuid:         summoner.Puuid,
		AccountId:     summoner.AccountId,
		SummonerId:    summoner.SummonerId,
		Level:         int64(summoner.Level),
		ProfileIconId: int32(summoner.ProfileIconId),
		Name:          summoner.Name,
		LastRevision:  time.Unix(int64(summoner.LastRevision), 0),
		TimeStamp:     time.Unix(summoner.TimeStamp, 0),
	}
}
