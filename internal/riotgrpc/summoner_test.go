package riotgrpc_test

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/rank1zen/yujin/internal/riotgrpc"
	"github.com/rank1zen/yujin/internal/riotgrpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func TestMain(m *testing.M) {
	lis = bufconn.Listen(bufSize)
	server := grpc.NewServer()

	golioClient := golio.NewClient(
		"RGAPI-b32ffbca-42f1-4ff1-8afe-241ab41fbbcc",
		golio.WithRegion(api.RegionNorthAmerica),
	)

	proto.RegisterRiotSummonerServer(server, riotgrpc.NewSummonerRpcServer(golioClient))
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("server exited with error: %v", err)
		}
	}()

	code := m.Run()

	defer server.Stop()
	defer lis.Close()

	os.Exit(code)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestOrrange(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewRiotSummonerClient(conn)

	expected := proto.Summoner{
		Puuid:      "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
		SummonerId: "2xCyr5bJbp2BlMSWLRolf9_x0eSbWBay5Bam_9myXFXjZSw",
		Name:       "orrange",
	}

	t.Run("ByNameQuery", func(t *testing.T) {
		summoner, err := client.ByName(
			ctx,
			&proto.ByNameRequest{Name: expected.Name},
		)
		if err != nil {
			t.Fatalf("error in querying summoner: %v", err)
		}
		if summoner.Puuid != expected.Puuid {
			t.Fatalf("wrong puuid")
		}
	})

	t.Run("GetMatchlistQuery", func(t *testing.T) {
		_, err := client.GetMatchlist(
			ctx,
			&proto.ByPuuidMatchlistRequest{
				Puuid: expected.Puuid,
				Start: 0,
				Count: 5,
			},
		)
		if err != nil {
			t.Fatalf("error in querying matchlist: %v", err)
		}
	})

	t.Run("GetSoloqQuery", func(t *testing.T) {
		_, err := client.GetSoloq(
			ctx,
			&proto.BySummonerIdRequest{SummonerId: expected.SummonerId},
		)
		if err != nil {
			t.Fatalf("error in querying soloq: %v", err)
		}
	})
}
