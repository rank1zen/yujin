package riotgrpc_test

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/rank1zen/yujin/internal/riotgrpc"
	"github.com/rank1zen/yujin/internal/riotgrpc/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestConnectWithInsecure(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.DialContext(ctx, os.Getenv("YUJIN_TEST_RPC_SERVER_ADDR"), opts...)
	require.NoError(t, err)

	closeConn(t, conn)
}

func TestByNameQuery(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.DialContext(ctx, os.Getenv("YUJIN_TEST_RPC_SERVER_ADDR"), opts...)
	require.NoError(t, err)
	defer closeConn(t, conn)

	client := proto.NewRiotSummonerClient(conn)
	_, err = client.ByName(ctx, &proto.ByNameRequest{Name: "orrange"})
	require.NoError(t, err)
}

func TestByNameQueryNotFoundError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.DialContext(ctx, os.Getenv("YUJIN_TEST_RPC_SERVER_ADDR"), opts...)
	require.NoError(t, err)
	defer closeConn(t, conn)

	client := proto.NewRiotSummonerClient(conn)
	_, err = client.ByName(ctx, &proto.ByNameRequest{Name: "o"})
	require.Error(t, err)
}

func TestMutliQuery(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.DialContext(ctx, os.Getenv("YUJIN_TEST_RPC_SERVER_ADDR"), opts...)
	require.NoError(t, err)
	defer closeConn(t, conn)

	client := proto.NewRiotSummonerClient(conn)

	summ, err := client.ByName(ctx, &proto.ByNameRequest{Name: "orrange"})
	require.NoError(t, err)

	_, err = client.GetMatchlist(ctx, &proto.ByPuuidMatchlistRequest{Puuid: summ.Puuid, Start: 0, Count: 10})
	require.NoError(t, err)

	_, err = client.GetSoloq(ctx, &proto.BySummonerIdRequest{SummonerId: summ.SummonerId})
	require.NoError(t, err)
}

func closeConn(t testing.TB, conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {
		t.Fatalf("conn.Close unexpectedly failed: %v", err)
	}
}

func TestMain(m *testing.M) {
	riotKey := os.Getenv("YUJIN_RIOT_API_KEY")
	if riotKey == "" {
		log.Fatalf("Skipping due to missing env %v", "YUJIN_RIOT_API_KEY")
	}

	gcOpts := []golio.Option{
		golio.WithRegion(api.RegionNorthAmerica),
	}
	gc := golio.NewClient(riotKey, gcOpts...)
	
	gs := grpc.NewServer()
	proto.RegisterRiotSummonerServer(gs, riotgrpc.NewSummonerRpcServer(gc))

	addr := os.Getenv("YUJIN_TEST_RPC_SERVER_ADDR")
	if addr == "" {
		log.Fatalf("Skipping due to missing env %v", "YUJIN_TEST_RPC_SERVER_ADDR")
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	go func() {
		gs.Serve(lis)
	}()

	code := m.Run()
	os.Exit(code)
}
