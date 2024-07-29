package database

import (
	"context"
	"fmt"
	"strings"
)

type RiotPuuid string

func (id RiotPuuid) String() string {
	return string(id)
}

type RiotSummonerId string

func (id RiotSummonerId) String() string {
	return string(id)
}

type RiotName struct {
	GameName string
	TagLine  string
}

func ParseRiotName(s string) (RiotName, error) {
	parts := strings.SplitN(s, "-", 2)
	return RiotName{
		GameName: parts[0],
		TagLine: parts[1],
	}, nil
}

func (name RiotName) String() string {
	return fmt.Sprintf("%s#%s", name.GameName, name.TagLine)
}

type Identifiers struct {
	Puuid      RiotPuuid
	SummonerId RiotSummonerId
}

func (db *DB) GetAccount(ctx context.Context, name RiotName) (*Identifiers, error) {
	acc, err := db.riot.GetAccountByRiotId(ctx, name.GameName, name.TagLine)
	if err != nil {
		return nil, fmt.Errorf("could not get riot account: %w", err)
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
