package riot

import (
	"context"
	"time"

	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/rank1zen/yujin/internal"
)

type SummonerQ struct {
	q *lol.SummonerClient
}

func NewSummonerQ(q *lol.SummonerClient) *SummonerQ {
	return &SummonerQ{
		q: q,
	}
}

func (s *SummonerQ) ByPuuid(ctx context.Context, puuid string) (Summoner, error) {
	timeStamp := time.Now().Unix()
	r, err := s.q.GetByPUUID(puuid)
	if err != nil {
		return Summoner{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "q.GetByPUUID")
	}
	return riotCast(r, timeStamp), nil
}

func (s *SummonerQ) ByName(ctx context.Context, name string) (Summoner, error) {
	timeStamp := time.Now().Unix()
	r, err := s.q.GetByName(name)
	if err != nil {
		return Summoner{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "q.GetByName")
	}
	return riotCast(r, timeStamp), nil
}

type Summoner struct {
	Puuid         string
	AccountId     string
	SummonerId    string
	Level         int
	ProfileIconId int
	Name          string
	LastRevision  int
	TimeStamp     int64
}

func riotCast(summoner *lol.Summoner, timeStamp int64) Summoner {
	return Summoner{
		Puuid:         summoner.PUUID,
		AccountId:     summoner.AccountID,
		SummonerId:    summoner.ID,
		Level:         summoner.SummonerLevel,
		ProfileIconId: summoner.ProfileIconID,
		Name:          summoner.Name,
		LastRevision:  summoner.RevisionDate,
		TimeStamp:     timeStamp,
	}
}
