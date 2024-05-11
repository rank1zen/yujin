package database

import (
	"fmt"
	"strings"
	"time"
)

type RecordFilter struct {
	Field string
	Value any
}

type fil interface {
	Eq() (string, any)
}

type SQLer struct {
	strings.Builder
	args []any
}


func (s *SQLer) With(f RecordFilter) {
	// s.WriteString(fmt.Sprintf(" AND %s = $%d", f.Field, len(args)))
}

func build(q string, start int, filters ...RecordFilter) (string, []any) {
	var args []any
	for _, f := range filters {
		args = append(args, f.Value)
		q += fmt.Sprintf(" AND %s = $%d", f.Field, len(args)+start)
	}

	strings.ReplaceAll(q, "\n", " ")
	return q, args
}

func BuildFilters[T any](arr []T) []RecordFilter {
	return nil
}

type SummonerRecordFilter struct {
	Field   string
	Value   string
	DateMin time.Time
	DateMax time.Time

	SortOrder string
	Offset    uint64
	Limit     uint64
}

type SummonerRecordCount struct {
	Field   string
	Value   string
	DateMin time.Time
	DateMax time.Time
}

func buildSummmonerFilterQuery(f *SummonerRecordFilter) (string, []interface{}) {
	var args []interface{}
	q := `
                SELECT
                        record_id, record_date, puuid, account_id, id,
                        name, profile_icon_id, summoner_level, revision_date
                FROM
                        summoner_records
                WHERE 1=1
        `

	if f.Field != "" {
		args = append(args, f.Value)
		q += fmt.Sprintf(" AND %s = $%d", f.Field, len(args))
	}

	if !f.DateMin.IsZero() {
		args = append(args, f.DateMin)
		q += fmt.Sprintf(" AND %s >= $%d", "record_date", len(args))
	}

	if !f.DateMax.IsZero() {
		args = append(args, f.DateMax)
		q += fmt.Sprintf(" AND %s < $%d", "record_date", len(args))
	}

	q = strings.ReplaceAll(q, "\n", " ")

	return q, args
}

func buildSummonerCountQuery(f *SummonerRecordCount) (string, []interface{}) {
	var args []interface{}
	q := `
                SELECT COUNT(*)
                FROM summoner_records
                WHERE 1=1
        `

	if f.Field != "" {
		args = append(args, f.Value)
		q += fmt.Sprintf(" AND %s = $%d", f.Field, len(args))
	}

	if !f.DateMin.IsZero() {
		args = append(args, f.DateMin)
		q += fmt.Sprintf(" AND %s >= $%d", "record_date", len(args))
	}

	if !f.DateMax.IsZero() {
		args = append(args, f.DateMax)
		q += fmt.Sprintf(" AND %s < $%d", "record_date", len(args))
	}

	q = strings.ReplaceAll(q, "\n", " ")

	return q, args
}
