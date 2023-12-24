package internal

import (
	"context"

	"github.com/rank1zen/yujin/internal/riotgrpc/proto"
	"google.golang.org/grpc"
)

type summone struct {
	sc proto.RiotSummonerClient
}

func NewSummonerGrpcClient(cc grpc.ClientConnInterface) *summone {
	return &summone{sc: proto.NewRiotSummonerClient(cc)}
}

func (s *summone) QueryRiot(ctx context.Context, name string) (SummonerWithIds, error) {
	r, err := s.sc.ByName(ctx, &proto.ByNameRequest{Name: name})
	if err != nil {
		return SummonerWithIds{}, err
	}
	return SummonerWithIds{
		Puuid: r.Puuid,
	}, nil
}

