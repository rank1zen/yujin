package internal

import (
	"context"

	"github.com/rank1zen/yujin/internal/postgresql/db"
)

type summonerClient struct {
	q  *db.Queries
}

func NewSummonerClient(d db.DBTX) *summonerClient {
	return &summonerClient{q: db.New(d)}
}

func (s *summonerClient) Create(ctx context.Context, params SummonerWithIds) (string, error) {
	uuid, err := s.q.InsertSummoner(ctx, params.CastToDB())
	if err != nil {
		return "", WrapErrorf(err, ErrorCodeUnknown, "s.q.InsertSummoner")
	}
	return UUIDString(uuid), nil
}

func (s *summonerClient) Find(ctx context.Context, puuid string, limit int32, offset int32) ([]Summoner, error) {
	summoner, err := s.q.SelectRecordsForSummoner(ctx, db.SelectRecordsForSummonerParams{
		Puuid:  puuid,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return []Summoner{}, WrapErrorf(err, ErrorCodeUnknown, "s.db.SelectRecordsForSummoner")
	}

	res := make([]Summoner, len(summoner))
	for i, sum := range summoner {
		res[i] = CastFromDB(&sum)
	}

	return res, nil
}

func (s *summonerClient) Newest(ctx context.Context, puuid string) (Summoner, error) {
	summoner, err := s.q.SelectRecentRecordForSummoner(ctx, puuid)
	if err != nil {
		return Summoner{}, WrapErrorf(err, ErrorCodeUnknown, "s.q.SelectRecentRecordForSummoner")
	}
	return CastFromDB(&summoner), nil
}

func (s *summonerClient) Delete(ctx context.Context, id string) error {
	idVal, err := ParseUUID(id)
	if err != nil {
		return WrapErrorf(err, ErrorCodeUnknown, "ParseUUID")
	}

	err = s.q.DeleteSummoner(ctx, idVal)
	if err != nil {
		return WrapErrorf(err, ErrorCodeUnknown, "s.q.DeleteSummoner")
	}
	return nil
}
