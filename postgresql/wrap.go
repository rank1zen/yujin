package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rank1zen/yujin/postgresql/gen"
)

type SummonerRecord struct {
	RecordDate time.Time
}

func InsertSummonerRecord(ctx context.Context, q *gen.Queries, r *SummonerRecord) (string, error) {
	params := gen.InsertSummonerRecordParams{
		//ProfileIconID: 1,
		//RecordDate: arg.RecordDate,
	}

	id, err := q.InsertSummonerRecord(ctx, params)
	if err != nil {
		return "", err
	}

	return UUIDString(id), nil
}

func DeleteSummonerRecord(ctx context.Context, q *gen.Queries, id string) (string, time.Time, error) {
	uuid, err := newUUID(id)
	if err != nil {
		return "", time.Time{}, err
	}

	row, err := q.DeleteSummonerRecord(ctx, uuid)
	if err != nil {
		return "", time.Time{}, err
	}

	return "", row.RecordDate.Time, nil
}

type SoloqRecord struct {
	RecordDate time.Time
}

func InsertSoloqRecord(ctx context.Context, q *gen.Queries, r *SoloqRecord) (string, error) {
	arg := gen.InsertSoloqRecordParams{}

	id, err := q.InsertSoloqRecord(ctx, arg)
	if err != nil {
		return "", err
	}

	return UUIDString(id,), nil
}

func DeletSoloqRecord(ctx context.Context, q *gen.Queries, id string) (string, time.Time, error) {
	uuid, err := newUUID(id)
	if err != nil {
		return "", time.Time{}, err
	}

	row, err := q.DeleteSoloqRecord(ctx, uuid)
	if err != nil {
		return "", time.Time{}, err
	}

	return "", row.RecordDate.Time, nil
}

func CountSoloqRecordsById(ctx context.Context, q *gen.Queries, id string) (int64, error) {
	var pg pgtype.Text
	pg.Scan(id)
	count, err := q.CountSoloqRecordsById(ctx, pgtype.Text{})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func CountSummonerRecordsByPuuid(ctx context.Context, q *gen.Queries, puuid string) (int64, error) {
	count, err := q.CountSummonerRecordsByPuuid(ctx, pgtype.Text{})
	if err != nil {
		return 0, err
	}

	return count, nil
}
