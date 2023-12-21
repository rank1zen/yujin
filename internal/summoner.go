package internal

import (
	"context"

	"github.com/rank1zen/yujin/internal/postgresql/db"
	"github.com/rank1zen/yujin/internal/riotgrpc/proto"
)

type SummonerClient struct {
	q *db.Queries
	sc proto.SummonerQueryClient
}

func NewSummonerClient(d db.DBTX, sc proto.SummonerQueryClient) *SummonerClient {
	return &SummonerClient{
		q: db.New(d),
		sc: sc,
	}
}

func (s *SummonerClient) Create(ctx context.Context, params SummonerWithIds) (string, error) {
	uuid, err := s.q.InsertSummoner(ctx, params.CastToDB())
	if err != nil {
		return "", err
	}
	return UUIDString(uuid), nil
}

func (s *SummonerClient) Find(ctx context.Context, puuid string, limit int32, offset int32) ([]Summoner, error) {
	summoner, err := s.q.SelectRecordsForSummoner(ctx, db.SelectRecordsForSummonerParams{
		Puuid:  puuid,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return []Summoner{}, err
	}

	res := make([]Summoner, len(summoner))
	for i, sum := range summoner {
		res[i].Level = sum.Level
		res[i].ProfileIconId = sum.ProfileIconID
		res[i].Name = sum.Name
		res[i].LastRevision = sum.LastRevision.Time
		res[i].TimeStamp = sum.TimeStamp.Time
	}

	return res, nil
}

func (s *SummonerClient) Newest(ctx context.Context, puuid string) (Summoner, error) {
	summoner, err := s.q.SelectRecentRecordForSummoner(ctx, puuid)
	if err != nil {
		return Summoner{}, err
	}
	return Summoner{
		Level:         summoner.Level,
		ProfileIconId: summoner.ProfileIconID,
		Name:          summoner.Name,
		LastRevision:  summoner.LastRevision.Time,
		TimeStamp:     summoner.TimeStamp.Time,
	}, nil
}

func (s *SummonerClient) Delete(ctx context.Context, id string) error {
	idVal, err := ParseUUID(id)
	if err != nil {
		return err
	}

	err = s.q.DeleteSummoner(ctx, idVal)
	if err != nil {
		return err
	}
	return nil
}
