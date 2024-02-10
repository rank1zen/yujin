package postgresql_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/internal/postgresql"
	"github.com/stretchr/testify/assert"
)

func TestMatchTeamRecord(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	addr := postgresql.NewDockerResource(t)

	pool := MustConnectionPool(t, addr)
	MustMigrate(t, pool)
	db := postgresql.NewMatchV5Query(pool)

	tests := []struct {
		arg  postgresql.MatchTeamRecordArg
		want string
	}{
		{
			arg:  postgresql.MatchTeamRecordArg{},
			want: "",
		},
	}

	for _, tc := range tests {
		id, err := db.InsertMatchTeam(ctx, &tc.arg)
		if assert.NoError(t, err) {
			assert.Equal(t, tc.want, id)
			t.Log(id)
		}
	}
}

func MustConnectionPool(t testing.TB, url string) (*pgxpool.Pool) {
	pool, err := postgresql.NewBackoffPool(context.Background(), url)
	if err != nil {
		t.Fatalf("can not make connection pool: %v", err)
	}

	return pool
}

func MustMigrate(t testing.TB, pool *pgxpool.Pool) {
	err := postgresql.Migrate(context.Background(), pool)
	if err != nil {
		t.Fatalf("can not migrate: %v", err)
	}
}
