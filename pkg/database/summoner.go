package database

import (
	"context"
	"time"
)

type SummonerRecord struct {
	RecordId      string    `db:"record_id"`
	RecordDate    time.Time `db:"record_date"`
	Puuid         string    `db:"puuid"`
	AccountId     string    `db:"account_id"`
	SummonerId    string    `db:"summoner_id"`
	ProfileIconId int32     `db:"profile_icon_id"`
	SummonerLevel int32     `db:"summoner_level"`
	RevisionDate  time.Time `db:"revision_date"`
}

type SummonerProfile struct {
	Name          string
	Level         int
	ProfileIconID int
	Wins          int
	Losses        int
	Rank          *string
	Tier          *string
	LP            *int
}

// FIXME: A simple sorting thing for summoners
// Which summoner (puuid), from dates A to B, Sort by Date ASC or DESC
// FIXME: pagination yes or no?
type SummonerRecordFilteR struct {
	Puuid   string
	DateMin time.Time
	DateMax time.Time
	DateAsc bool
}

// TODO: implement
func (r *service) FetchAndInsert(ctx context.Context, puuid string) error {
	return nil
}

// TODO: implement
func (r *service) GetRecent(ctx context.Context, puuid string) (*SummonerRecord, error) {
	return nil, nil
}
