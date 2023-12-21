package riotgrpc_test

import (
	"context"
	"log"
	"net"
	_"os"
	"testing"

	"github.com/KnutZuidema/golio"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"github.com/rank1zen/yujin/internal/riotgrpc/proto"
)

const bufSize = 1024 * 1024
var lis *bufconn.Listener

func TestMain(m *testing.M) {
	lis = bufconn.Listen(bufSize)
	//grpcServer := grpc.NewServer()


	apiKey := "RGAPI-88780330-dc2d-4ff0-8605-95d03d2c22c9"
	log.Printf("API KEY: %s", apiKey)

	golioClient := golio.NewClient(apiKey)
	log.Printf("Ok golio up")

	champ, _ := golioClient.DataDragon.GetChampion("Ashe")
	log.Printf(champ.ID)

	//summonerServer := riot.NewSummonerClient(golioClient)
	//riot.RegisterSummonerQueryServer(grpcServer, summonerServer)

	//go func() {
		//if err := grpcServer.Serve(lis); err != nil {
			//log.Fatalf("server exited with error: %v", err)
		//}
	//}()

	//code := m.Run()
	//grpcServer.Stop()
	//os.Exit(code)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestBasicQuery(t *testing.T) {
	t.Logf("running client test")
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(bufDialer),
	)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := proto.NewSummonerQueryClient(conn)
	summoner, err := client.ByName(ctx, &proto.SummonerByNameRequest{Name:"orrange"})
	if err != nil {
		t.Fatalf("error in querying summoner: %v", err)
	}
	t.Logf("success: %v", summoner)
}
