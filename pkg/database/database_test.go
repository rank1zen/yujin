package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/pkg/docker"
	"github.com/rank1zen/yujin/pkg/logging"
	"github.com/stretchr/testify/assert"
)

var conn *pgx.Conn
var mu sync.Mutex

func TestMain(m *testing.M) {
	var purge func()
	conn, purge = docker.NewPostgres()
	defer purge()

	code := m.Run()
	os.Exit(code)
}

func setupDB(tb testing.TB) *DB {
	tb.Helper()

	ctx := context.Background()

	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		tb.Fatalf("bytes: %v", err)
	}

	randName := hex.EncodeToString(b)

	// postgres does not allow parallel database creation
	mu.Lock()
	defer mu.Unlock()
	_, err = conn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s" WITH TEMPLATE "%s";`, randName, conn.Config().Database))
	if err != nil {
		tb.Fatalf("failed to clone db: %v", err)
	}

	newCfg := conn.Config()
	newCfg.Database = randName

	db, err := NewDB(ctx, newCfg.ConnString())
	if err != nil {
		tb.Fatalf("failed to connect db: %v", err)
	}

	tb.Cleanup(func() {
		ctx := context.Background()

		db.Close()

		mu.Lock()
		defer mu.Unlock()
		_, err := conn.Exec(ctx, fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE);`, randName))
		if err != nil {
			tb.Errorf("failed to drop cloned db: %v", err)
		}
	})

	return db
}

func NewPool(tb testing.TB) *pgxpool.Pool {
	tb.Helper()

	ctx := context.Background()

	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		tb.Fatalf("bytes: %v", err)
	}

	randName := hex.EncodeToString(b)

	// postgres does not allow parallel database creation
	mu.Lock()
	defer mu.Unlock()
	_, err = conn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s" WITH TEMPLATE "%s";`, randName, conn.Config().Database))
	if err != nil {
		tb.Fatalf("failed to clone db: %v", err)
	}

	newCfg := conn.Config()
	newCfg.Database = randName

	poolCfg, err := pgxpool.ParseConfig(newCfg.ConnString())
	if err != nil {
		tb.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		tb.Fatalf("failed to connect db: %s", err)
	}

	tb.Cleanup(func() {
		ctx := context.Background()

		pool.Close()

		mu.Lock()
		defer mu.Unlock()

		_, err := conn.Exec(ctx, fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE);`, conn.Config().Database))
		if err != nil {
			tb.Errorf("failed to drop cloned db: %v", err)
		}
	})

	return pool
}

func testingContext(tb testing.TB) context.Context {
	ctx := context.Background()
	return logging.WithContext(ctx, logging.NewTestLogger(tb))
}

func TestDockerConn(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	db := NewPool(t)

	err := db.Ping(ctx)
	assert.NoError(t, err)
}
