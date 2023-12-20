package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"github.com/rank1zen/yujin/internal/riot"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := *grpc.NewServer(opts...)
	s := &riot.Summoner{}
	riot.RegisterSummonerQueryServer(&grpcServer, s)
	log.Print("I guess the server is running:?")
	grpcServer.Serve(lis)
}
