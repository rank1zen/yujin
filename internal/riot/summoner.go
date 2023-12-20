package riot

import (
	"context"

	"github.com/KnutZuidema/golio"
)

type Summoner struct {
	UnimplementedSummonerQueryServer
	gc *golio.Client
}

func NewSummonerClient(gc *golio.Client) *Summoner {
	return &Summoner{gc: gc}
}

func (s *Summoner) ByName(ctx context.Context, req *SummonerByNameRequest) (*SummonerResponse, error) {
	lolSum, err := s.gc.Riot.LoL.Summoner.GetByName(req.GetName())
	if err != nil {
		return &SummonerResponse{}, err
	}
	return &SummonerResponse{
		Puuid:         lolSum.PUUID,
		AccountId:     lolSum.AccountID,
		SummonerId:    lolSum.ID,
		Level:         int64(lolSum.SummonerLevel),
		ProfileIconId: int32(lolSum.ProfileIconID),
		Name:          lolSum.Name,
		LastRevision:  int64(lolSum.RevisionDate),
	}, nil
}

func (s *Summoner) ByPuuid(ctx context.Context, req *SummonerByPuuidRequest) (*SummonerResponse, error) {
	return &SummonerResponse{}, nil
}
