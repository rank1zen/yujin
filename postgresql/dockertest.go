package postgresql

import (
	"fmt"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func NewDockerResource(t testing.TB) string {
	t.Log("constructing dockertest pool")
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("could not construct dockertest pool: %v", err)
	}

	opts :=  dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15-alpine3.18",
		Env: []string{
			"POSTGRES_PASSWORD=yuyu",
			"listen_addresses = '*'",
		},
	}

	hc := func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
	}

	t.Log("starting dockertest resource")
	resource, err := pool.RunWithOptions(&opts, hc)
	if err != nil {
		t.Fatalf("could not start dockertest resource: %v", err)
	}

	resource.Expire(60)

	t.Cleanup(func() {
		t.Log("puring dockertest resource")
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("could'nt purge container: %v", err)
		}
	})

	return fmt.Sprintf("postgres://postgres:yuyu@%s/postgres?sslmode=disable", resource.GetHostPort("5432/tcp"))
}
