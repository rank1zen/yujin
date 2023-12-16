package postgresql

import (
	"context"

	"github.com/google/uuid"
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

func (s *SummonerDA) Create(ctx context.Context, summoner internal.Summoner) (string, error) {
	id, err := s.q.InsertSummoner(ctx, db.InsertSummonerParams{

	})
	if err != nil {
		return "", err
	}
	return *id, nil
}

func (s *SummonerDA) Find(ctx context.Context, puuid string) ([]internal.Summoner, error) {
	summoner, err := s.q.SelectRecordsForSummoner(ctx, db.SelectRecordsForSummonerParams{

	})
	if err != nil {
		return nil, err
	}

	res := make([]internal.Summoner, len(summoner))
	for i, sum := range summoner {
		res[i] = dbCast(sum)
	}

	
	return res, nil
}

func (s *SummonerDA) Newest(ctx context.Context, puuid string) (internal.Summoner, error) {
	summoner, err := s.q.SelectRecentRecordForSummoner(ctx, puuid)
	if err != nil {
		return internal.Summoner{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "q.SelectRecentRecord...")
	}
	return dbCast(summoner), nil
}

func (s * SummonerDA) Delete(ctx context.Context, id string) error {
	val, err := uuid.Parse(id)

	err := s.q.DeleteSummoner(ctx, val)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "q.DeleteSummoner")
	}
	return nil
}

func dbCast(s db.Summoner) internal.Summoner {
	return internal.Summoner{

	}
}
