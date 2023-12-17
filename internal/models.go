package internal

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Summoner struct {
	Level         int64
	ProfileIconId int32
	Name          string
	LastRevision  pgtype.Timestamp
	TimeStamp     pgtype.Timestamp
}

type SummonerParams struct {
	Puuid         string
	AccountId     string
	SummonerId    string
	Level         int64
	ProfileIconId int32
	Name          string
	LastRevision  pgtype.Timestamp
	TimeStamp     pgtype.Timestamp
}
