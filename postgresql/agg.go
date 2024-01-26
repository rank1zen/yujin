package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rank1zen/yujin/postgresql/gen"
)

func CheckConn(pool *pgxpool.Pool) error {
	return pool.Ping(context.Background())
}

func NewPool(ctx context.Context, url string, e *echo.Echo) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool

	op := func() error {
		pool, err := pgxpool.New(ctx, url)
		if err != nil {
			return err
		}
		return pool.Ping(ctx)
	}

	b := backoff.NewExponentialBackOff()

	if err := backoff.RetryNotify(op, backoff.WithMaxRetries(b, 4)), func(error, time.Duration) {  }); err != nil {
		return nil, err
	}

	return pool, nil
}

type SummonerRecord struct {
	RecordDate time.Time
}

func InsertSummonerRecord(q *gen.Queries, ctx context.Context, arg *SummonerRecord) (string, error) {
	params := gen.InsertSummonerRecordParams{
		//RecordDate: arg.RecordDate,
	}

	id, err := q.InsertSummonerRecord(ctx, params)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", id.Bytes[0:4], id.Bytes[4:6], id.Bytes[6:8], id.Bytes[8:10], id.Bytes[10:16]), nil
}

func DeleteSummonerRecord(ctx context.Context, q *gen.Queries, id string) (string, time.Time, error) {
	uuid, err := ParseUUID(id)
	if err != nil {
		return "", time.Time{}, err
	}

	row, err := q.DeleteSummonerRecord(ctx, uuid)
	if err != nil {
		return "", time.Time{}, err
	}

	return "", row.RecordDate.Time, nil
}
