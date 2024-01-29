package postgresql

import (
	"fmt"
	"testing"

	"github.com/ory/dockertest/v3"
)

func NewDockerResource(t testing.TB) string {
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
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("failed to purge postgres container: %s", err)
		}
		t.Log("OK: purged postgres container")
	})

	return fmt.Sprintf("postgres://postgres:yuyu@%s/postgres?sslmode=disable", resource.GetHostPort("5432/tcp"))
}
