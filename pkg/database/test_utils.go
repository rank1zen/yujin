package database

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	repository       = "postgres"
	tag              = "alpine"
	databaseName     = "test_db"
	databaseUser     = "test_user"
	databasePassword = "testing123"
)

type TestInstance struct {
	SkipDB     bool
	SkipReason string

	TestDB *TestDatabaseResource
}

// An instance of a Docker container to run DB tests against
type TestDatabaseResource struct {
	pool      *dockertest.Pool
	container *dockertest.Resource
	Url       *url.URL
}

func NewTestInstance() *TestInstance {
	if !flag.Parsed() {
		flag.Parse()
	}

	if testing.Short() {
		return &TestInstance{
			SkipDB:     true,
			SkipReason: "-short flag provided",
			TestDB:     nil,
		}
	}
	return &TestInstance{
		SkipDB:     false,
		SkipReason: "",
		TestDB:     MustDatabaseResource(),
	}
}

func MustDatabaseResource() *TestDatabaseResource {
	db, err := NewDatabaseResource()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	return db
}

func NewDatabaseResource() (*TestDatabaseResource, error) {
	log.Printf("Connecting to Docker")
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Docker: %w", err)
	}

	log.Printf("Running Docker Container: %s:%s", repository, tag)
	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: repository,
		Tag:        tag,
		Env: []string{
			"POSTGRES_DB=" + databaseName,
			"POSTGRES_USER=" + databaseUser,
			"POSTGRES_PASSWORD=" + databasePassword,
		},
	},
		func(c *docker.HostConfig) {
			c.AutoRemove = true
			c.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	container.Expire(120)

	connUrl := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(databaseUser, databasePassword),
		Host:     container.GetHostPort("5432/tcp"),
		Path:     databaseName,
		RawQuery: "sslmode=disable",
	}

        time.Sleep(5 * time.Second)

	return &TestDatabaseResource{
		pool:      pool,
		container: container,
		Url:       connUrl,
	}, nil
}

func (t *TestDatabaseResource) MustClose() {
	err := t.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func (t *TestDatabaseResource) Close() error {
	err := t.pool.Purge(t.container)
	if err != nil {
		return err
	}

	return nil
}

func (t *TestDatabaseResource) NewDB(tb testing.TB) *DB {
	tb.Helper()

	var connstr string

	pool, err := pgxpool.New(context.Background(), connstr)
	if err != nil {
		tb.Skipf("could not connect to DB")
	}

	db := &DB{Pgx: pool}

	tb.Cleanup(func() {
		ctx := context.Background()

		db.Close(ctx)
	})

	return db
}
