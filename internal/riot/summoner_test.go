package riot_test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/rank1zen/yujin/internal/riot"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	port       = flag.Int("port", 50051, "The test grpc server port")
)

func TestMain(m *testing.M) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	golioClient := golio.NewClient(os.Getenv("RIOT_API_KEY"), golio.WithRegion(api.RegionNorthAmerica), golio.WithLogger(logrus.New()))
	s := riot.NewSummonerClient(golioClient)

	riot.RegisterSummonerQueryServer(grpcServer, s)

	grpcServer.Serve(lis)
}

func TestBasicQuery(t *testing.T) {
	conn, err := grpc.Dial(*serverAddr)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	client := riot.NewSummonerQueryClient(conn)
	defer conn.Close()

	ctx, cFunc := context.WithTimeout(context.Background(), 10)
	defer cFunc()
	summoner, err := client.ByName(ctx, &riot.SummonerByNameRequest{Name:"orrange"})
	if err != nil {
		t.Fatalf("error in querying summoner: %v", err)
	}

	t.Logf("success: %v", summoner)
}
