package postgresql

import (
	"context"
	"time"
)

func InsertSoloqRecord(ctx context.Context, r *SoloqRecord) (string, error) {
	return "", nil
}

func CountSoloqRecordsById(ctx context.Context, id string) (int64, error) {
	return 0, nil
}

func DeletSoloqRecord(ctx context.Context, id string) (string, time.Time, error) {
	return "", time.Time{}, nil
}
