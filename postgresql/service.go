package postgresql

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/postgresql/db"
)

func Insert(pool *pgxpool.Pool) func() error {
	return func() error {
		pool.Begin()
	}
}
