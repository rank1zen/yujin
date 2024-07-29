package database

import (
	"context"
	"time"
)

type SummonerRecord struct {
	RecordDate    time.Time `db:"record_date"`
	RevisionDate  time.Time `db:"revision_date"`
	RecordId      string    `db:"record_id"`
	Puuid         string    `db:"puuid"`
	AccountId     string    `db:"account_id"`
	SummonerId    string    `db:"summoner_id"`
	ProfileIconId int32     `db:"profile_icon_id"`
	SummonerLevel int32     `db:"summoner_level"`
}

func (db *DB) updateSummoner(ctx context.Context, puuid RiotPuuid) error {
	m, err := db.riot.GetSummoner(ctx, puuid.String())
	if err != nil {
		return err
	}

	row := map[string]any{}
	row["summoner_id"] = m.Id
	row["account_id"] = m.AccountId
	row["puuid"] = m.Puuid
	row["revision_date"] = time.Unix(m.RevisionDate/1000, 0)
	row["profile_icon_id"] = m.ProfileIconId
	row["summoner_level"] = m.SummonerLevel

	// TODO: make these go in transaction
	err = queryInsertRow(ctx, db.pool, "summoner_records", row)
	if err != nil {
		return err
	}

	return nil
}
