package internal

import "time"

type SummonerRecord struct {
	ID            SummonerID
	Puuid         PUUID
	RevisionDate  time.Time
	ValidFrom     time.Time
	ValidTo       time.Time
	EnteredAt     time.Time
	Level         int
	ProfileIconID ProfileIconID
}

type RankRecord struct {
	Puuid     PUUID
	Timestamp time.Time
	LeagueID string
	Wins      int
	Losses    int
	Tier      int
	Division  string
	LP        int
}

type RankSnapshot struct {
	RankRecord

	LpDelta *int
}
