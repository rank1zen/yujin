package database

import "time"

type SummonerRecord struct {
	RecordId      string    `db:"record_id"`
	RecordDate    time.Time `db:"record_date"`
	Puuid         string    `db:"puuid"`
	AccountId     string    `db:"account_id"`
	SummonerId    string    `db:"id"`
	Name          string    `db:"name"`
	ProfileIconId int32     `db:"profile_icon_id"`
	SummonerLevel int32     `db:"summoner_level"`
	RevisionDate  int64     `db:"revision_date"`
}
