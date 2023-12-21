package internal_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rank1zen/yujin/internal"
)

var dbpool *pgxpool.Pool

func TestMain(m *testing.M) {
	dockerPool := initDocker()
	resource := createContainer(dockerPool)

	databaseUrl := fmt.Sprintf("postgres://yuyu:yuyu@localhost:%s/summoners?sslmode=disable", resource.GetPort("5432/tcp"))

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := dockerPool.Retry(
		func() error {
			ctx := context.Background()
			var err error
			dbpool, err = pgxpool.New(ctx, databaseUrl)
			if err != nil {
				return err
			}
			return dbpool.Ping(ctx)
		},
	); err != nil {
		log.Fatalf("Could not pool to database: %s", err)
	}

	ctx := context.Background()

	db, err := pgx.Connect(ctx, databaseUrl)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	migrator, err := migrate.NewMigrator(ctx, db, "public.schema_version")
	if err != nil {
		log.Fatalf("Could not create migrator: %s", err)
	}

	if err = migrator.LoadMigrations(os.DirFS("../db/migrations")); err != nil {
		log.Fatalf("Could not load migrations: %s", err)
	}

	if err = migrator.Migrate(ctx); err != nil {
		log.Fatalf("Could not migrate: %s", err)
	}

	code := m.Run()

	if err := dockerPool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func initDocker() *dockertest.Pool {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	pool.MaxWait = 60 * time.Second
	return pool
}

func createContainer(pool *dockertest.Pool) *dockertest.Resource {
	// pulls an image, creates a container based on it and runs it
	container, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "15-alpine3.18",
			Env: []string{
				"POSTGRES_PASSWORD=yuyu",
				"POSTGRES_USER=yuyu",
				"POSTGRES_DB=summoners",
				"listen_addresses = '*'",
			},
		},
		func(hc *docker.HostConfig) {
			hc.AutoRemove = true
			hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	container.Expire(60)
	return container
}

func TestSummonerCreateAndDelete(t *testing.T) {
	q := internal.NewSummonerClient(dbpool)
	id, err := q.Create(
		context.Background(),
		internal.SummonerWithIds{
			Puuid:         "YUYU",
			AccountId:     "YUYU",
			SummonerId:    "YUYU",
			Level:         324,
			ProfileIconId: 1008,
			Name:          "YUYU",
			LastRevision:  time.Now(),
			TimeStamp:     time.Now(),
		},
	)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	t.Logf("created with ID: %s", id)

	err = q.Delete(context.Background(), id)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
}
