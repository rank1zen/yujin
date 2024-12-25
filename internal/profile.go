package internal

import "time"

// Profile represents a record of a summoner.
type Profile struct {
	RecordDate    time.Time
	RevisionDate  time.Time
	Puuid         PUUID
	SummonerID    SummonerID
	AccountID     AccountID
	Name          string
	Tagline       string
	Level         int
	ProfileIconID ProfileIconID
	Rank          *RankRecord
}

type ProfileHeader struct {
	Puuid       PUUID
	LastUpdated time.Time
	RiotID      string
	RiotTagLine string
	Rank        *RankRecord
}
