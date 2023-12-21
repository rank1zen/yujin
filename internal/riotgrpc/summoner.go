package riotgrpc

import (
	"context"

	"github.com/KnutZuidema/golio"
	pb "github.com/rank1zen/yujin/internal/riotgrpc/proto"
)

type SummonerServer struct {
	pb.UnimplementedSummonerQueryServer
	gc *golio.Client
}

func NewSummonerServer(gc *golio.Client) *SummonerServer {
	return &SummonerServer{gc: gc}
}

func (s *SummonerServer) ByName(ctx context.Context, req *pb.SummonerByNameRequest) (*pb.SummonerResponse, error) {
	lolSum, err := s.gc.Riot.LoL.Summoner.GetByName(req.GetName())
	if err != nil {
		return &pb.SummonerResponse{}, err
	}
	return &pb.SummonerResponse{
		Puuid:         lolSum.PUUID,
		AccountId:     lolSum.AccountID,
		SummonerId:    lolSum.ID,
		Level:         int64(lolSum.SummonerLevel),
		ProfileIconId: int32(lolSum.ProfileIconID),
		Name:          lolSum.Name,
		LastRevision:  int64(lolSum.RevisionDate),
	}, nil
}

func (s *SummonerServer) ByPuuid(ctx context.Context, req *pb.SummonerByPuuidRequest) (*pb.SummonerResponse, error) {
	return &pb.SummonerResponse{}, nil
}
