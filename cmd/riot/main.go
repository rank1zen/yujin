package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/rank1zen/yujin/internal/riotgrpc"
	"github.com/rank1zen/yujin/internal/riotgrpc/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	key := os.Getenv("YUJIN_RIOT_API_KEY")
	if key == "" {
		log.Fatal("no riot key")
	}

	gcOpts := []golio.Option{
		golio.WithRegion(api.RegionNorthAmerica),
		golio.WithClient(http.DefaultClient),
		golio.WithLogger(logrus.StandardLogger()),
	}

	gc := golio.NewClient(key, gcOpts...)
	
	gs := grpc.NewServer()
	proto.RegisterRiotSummonerServer(gs, riotgrpc.NewSummonerRpcServer(gc))

	lis, err := net.Listen("tcp", "127.0.0.1:5010")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	log.Printf("listening on: %s", lis.Addr().String())

	gs.Serve(lis)
}
