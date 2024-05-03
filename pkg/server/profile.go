package server

import (
	"context"

	"github.com/rank1zen/yujin/pkg/proto"
)

type summonerProfile struct {
        proto.UnimplementedSummonerProfileServer
}

func (s *summonerProfile) GetSummoner(ctx context.Context, req *proto.GetSummonerRequest) (*proto.GetSummonerResponse, error) {
        return &proto.GetSummonerResponse{
                Summoner: &proto.SummonerRecord{
                        SummonerName: "TestTest",
                },
        }, nil
}
