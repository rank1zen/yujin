package database

import (
	"context"
	"reflect"

	"github.com/jackc/pgx/v5"
)

func newSummonerRecord[M any](v ...M) ([]SummonerRecord, error) {
        return nil, nil
}

func NewMatchBanRecord() []MatchBanRecord {
        return nil
}

func NewMatchObjectiveRecord() []MatchObjectiveRecord {
        return nil
}

func NewLeagueRecord() []LeagueRecord {
        return nil
}

// NewMatchRecord looks for struct tag json in M and returns an array of MatchRecord
func NewMatchRecord[M any](v ...M) []MatchRecord {
        for _, vv := range v {
                t := reflect.TypeOf(vv)
                for i := 0; i < t.NumField(); i++ {
                        field := t.Field(i)
                }
        }

        return nil
}

func NewMatchTeamRecord[M any](v ...M) []MatchTeamRecord {
        return nil
}

// extractStructSlice returns the field `db` columns and the values as 2d arrays
// Can take a slice of struct or a slice of pointers to structs
func extractStructSlice[T any](a []T) ([]string, [][]any, error) {
	var cols []string
        rows := make([][]any, len(a))

	t := reflect.TypeOf(a)

	t = t.Elem()

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("ok")
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		cols = append(cols, f.Tag.Get("db"))
	}

        row := make([]any, t.NumField())
        for i := range len(a) {
                for j := 0; j < t.NumField(); j++ {
                        row[j] = reflect.ValueOf(a[i]).FieldByName([]int{j})
                }
                rows[i] = row
        }

	return cols, rows, nil
}

func insertBulk[T any](ctx context.Context, db pgxDB, table string, records []T) (int64, error) {
	fields, rows, err := extractStructSlice(records)
	count, err := db.CopyFrom(ctx, pgx.Identifier{table}, fields, pgx.CopyFromRows(rows))
	if err != nil {
		return 0, err
	}

	return count, nil
}
