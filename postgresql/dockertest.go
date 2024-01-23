package postgresql

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
)

func NewConnection(t testing.TB) *pgx.Conn {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	log.Println("Running docker container (POSTGRES 15-Alpine3.18)")
	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "15-alpine3.18",
			Env: []string{
				"POSTGRES_PASSWORD=yuyu",
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

	resource.Expire(60)
	t.Cleanup(func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Couldn't purge container :%v", err)
		}
	})

	databaseUrl := fmt.Sprintf("postgres://postgres:yuyu@%s/postgres?sslmode=disable", resource.GetHostPort("5432/tcp"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	log.Println("Connecting to database on url: ", databaseUrl)
	var conn *pgx.Conn
	err = pool.Retry(func() error {
		var err error
		conn, err = pgx.Connect(ctx, databaseUrl)
		if err != nil {
			return err
		}
		return conn.Ping(ctx)
	})
	if err != nil {
		log.Fatalf("Could not pool to database: %s", err)
	}
	log.Println("Ok Connected to database")
	t.Cleanup(func(){
		conn.Close(ctx)
	})

	return conn
}

func MigrateDB(ctx context.Context, conn *pgx.Conn) {
	migrator, err := migrate.NewMigrator(ctx, conn, "public.schema_version")
	if err != nil {
		log.Fatalf("Could not create migrator: %s", err)
	}

	err = migrator.LoadMigrations(os.DirFS("../db/migrations"))
	if err != nil {
		log.Fatalf("Could not load migrations: %s", err)
	}

	err = migrator.Migrate(ctx)
	if err != nil {
		log.Fatalf("Could not migrate: %s", err)
	}
}
