package riotgrpc

import (
	"context"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/rank1zen/yujin/internal/riotgrpc/proto"
)

const SOLO_QUEUE = "RANKED_SOLO_5x5"

var queue = 420

type Summoner struct {
	proto.UnimplementedRiotSummonerServer
	gc *golio.Client
}

func NewSummonerRpcServer(gc *golio.Client) *Summoner {
	return &Summoner{gc: gc}
}

func (s *Summoner) ByName(ctx context.Context, req *proto.ByNameRequest) (*proto.Summoner, error) {
	name := req.GetName()
	ls, err := s.gc.Riot.LoL.Summoner.GetByName(name)
	if err != nil {
		return &proto.Summoner{}, err
	}
	return cast(ls), nil
}

func (s *Summoner) ByPuuid(ctx context.Context, req *proto.ByPuuidRequest) (*proto.Summoner, error) {
	puuid := req.GetPuuid()
	ls, err := s.gc.Riot.LoL.Summoner.GetByPUUID(puuid)
	if err != nil {
		return &proto.Summoner{}, err
	}
	return cast(ls), nil
}

func (s *Summoner) BySummonerId(ctx context.Context, req *proto.BySummonerIdRequest) (*proto.Summoner, error) {
	sumId := req.GetSummonerId()
	ls, err := s.gc.Riot.LoL.Summoner.GetByID(sumId)
	if err != nil {
		return &proto.Summoner{}, err
	}
	return cast(ls), nil
}

func (s *Summoner) GetSoloq(ctx context.Context, req *proto.BySummonerIdRequest) (*proto.LeagueEntry, error) {
	summonerId := req.GetSummonerId()
	entries, err := s.gc.Riot.LoL.League.ListBySummoner(summonerId)
	if err != nil {
		return &proto.LeagueEntry{}, err
	}

	soloq := &proto.LeagueEntry{SummonerId: summonerId}
	for _, entry := range entries {
		if entry.QueueType == SOLO_QUEUE {
			soloq.SummonerName = entry.SummonerName
			soloq.Tier = entry.Tier
			soloq.Rank = entry.Rank
			soloq.Lp = int32(entry.LeaguePoints)
			soloq.Wins = int32(entry.Wins)
			soloq.Losses = int32(entry.Losses)
		}
	}

	return soloq, nil
}

func (s *Summoner) GetMatchlist(ctx context.Context, req *proto.ByPuuidMatchlistRequest) (*proto.Matchlist, error) {
	ml, err := s.gc.Riot.LoL.Match.List(
		req.GetPuuid(),
		int(req.GetStart()),
		int(req.GetCount()),
		&lol.MatchListOptions{Queue: &queue},
	)
	if err != nil {
		return &proto.Matchlist{}, err
	}
	return &proto.Matchlist{MatchIds: ml}, nil
}

func cast(s *lol.Summoner) *proto.Summoner {
	return &proto.Summoner{
		Puuid:         s.PUUID,
		AccountId:     s.AccountID,
		SummonerId:    s.ID,
		Level:         int64(s.SummonerLevel),
		ProfileIconId: int32(s.ProfileIconID),
		Name:          s.Name,
		LastRevision:  int64(s.RevisionDate),
	}
}
