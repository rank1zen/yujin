package database

import (
	"context"
	"fmt"
)

type RiotPuuid string

type RiotSummonerId string

type RiotName struct {
	GameName string
	TagLine  string
}

func ParseRiotName(s string) (RiotName, error) {
	return RiotName{
		GameName: "a",
		TagLine: "a",
	}, nil
}

func (n RiotName) String() string {
	return fmt.Sprintf("%s #%s", n.GameName, n.TagLine)
}

type Identifiers struct {
	Puuid      RiotPuuid
	SummonerId RiotSummonerId
}

func (db *DB) GetAccount(ctx context.Context, name RiotName) (*Identifiers, error) {
	acc, err := db.riot.GetAccountByRiotId(ctx, name.GameName, name.TagLine)
	if err != nil {
		return nil, err
	}

	summ, err := db.riot.GetSummoner(ctx, acc.Puuid)
	if err != nil {
		return nil, err
	}

	return &Identifiers{
		Puuid:      RiotPuuid(acc.Puuid),
		SummonerId: RiotSummonerId(summ.Id),
	}, nil
}

func (db *DB) CheckFirstTimeVisit(ctx context.Context, puuid string) (bool, error) {
	var found bool
	err := db.pool.QueryRow(ctx, `
	SELECT EXISTS (SELECT 1 FROM summoner_records WHERE puuid = $1)
	`, puuid).Scan(&found)

	return found, err
}
