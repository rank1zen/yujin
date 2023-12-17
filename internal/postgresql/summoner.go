package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/postgresql/db"
)

type SummonerDA struct {
	q *db.Queries
}

func NewSummonerDA(d db.DBTX) *SummonerDA {
	return &SummonerDA{q: db.New(d)}
}

func (s *SummonerDA) Create(ctx context.Context, params internal.SummonerParams) (pgtype.UUID, error) {
	id, err := s.q.InsertSummoner(ctx, db.InsertSummonerParams{
		Puuid:         params.Puuid,
		AccountID:     params.AccountId,
		SummonerID:    params.SummonerId,
		Level:         params.Level,
		ProfileIconID: params.ProfileIconId,
		Name:          params.Name,
		LastRevision:  params.LastRevision,
		TimeStamp:     params.TimeStamp,
	})
	if err != nil {
		return pgtype.UUID{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "q.InsertSummoner")
	}
	return id, nil
}

func (s *SummonerDA) Find(ctx context.Context, puuid string, limit int32, offset int32) ([]internal.Summoner, error) {
	summoner, err := s.q.SelectRecordsForSummoner(ctx, db.SelectRecordsForSummonerParams{
		Puuid:  puuid,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return []internal.Summoner{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "q.SelectRecords...")
	}

	res := make([]internal.Summoner, len(summoner))
	for i, sum := range summoner {
		res[i].Level = sum.Level
		res[i].ProfileIconId = sum.ProfileIconID
		res[i].Name = sum.Name
		res[i].LastRevision = sum.LastRevision
		res[i].TimeStamp = sum.TimeStamp
	}

	return res, nil
}

func (s *SummonerDA) Newest(ctx context.Context, puuid string) (internal.Summoner, error) {
	summoner, err := s.q.SelectRecentRecordForSummoner(ctx, puuid)
	if err != nil {
		return internal.Summoner{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "q.SelectRecentRecord...")
	}
	return internal.Summoner{
		Level:         summoner.Level,
		ProfileIconId: summoner.ProfileIconID,
		Name:          summoner.Name,
		LastRevision:  summoner.LastRevision,
		TimeStamp:     summoner.TimeStamp,
	}, nil
}

func (s *SummonerDA) Delete(ctx context.Context, id pgtype.UUID) error {
	err := s.q.DeleteSummoner(ctx, id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "q.DeleteSummoner")
	}
	return nil
}
