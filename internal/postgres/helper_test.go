package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/rank1zen/yujin/internal/postgres"
)

type ConnTestRunner struct {
	CreateConfig func(ctx context.Context, t testing.TB) *pgx.ConnConfig
	AfterConnect func(ctx context.Context, t testing.TB, conn *pgx.Conn)
	AfterTest    func(ctx context.Context, t testing.TB, conn *pgx.Conn)
	CloseConn    func(ctx context.Context, t testing.TB, conn *pgx.Conn)
}

func NewConnTestRunner() *ConnTestRunner {
	return &ConnTestRunner{}
}

func NewDockerResource(t testing.TB) string {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("could not construct dockertest pool: %v", err)
	}

	t.Log("OK: constructed dockertest pool")

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16",
		Env: []string{
			"POSTGRES_PASSWORD=yuyu",
			"listen_addresses='*'",
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		t.Fatalf("could not start dockertest resource: %v", err)
	}

	t.Log("OK: docker resource up")

	resource.Expire(60)

	t.Cleanup(func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("failed to purge postgres container: %s", err)
		}

		t.Log("OK: purged postgres container")
	})

	return fmt.Sprintf("postgres://postgres:yuyu@%s/postgres?sslmode=disable", resource.GetHostPort("5432/tcp"))
}

func MustConnectionPool(t testing.TB, url string) *pgxpool.Pool {
	pool, err := postgres.NewBackoffPool(context.Background(), url)
	if err != nil {
		t.Fatalf("can not make connection pool: %v", err)
	}

	return pool
}
