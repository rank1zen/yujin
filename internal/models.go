package internal

import (
	"time"
)

type Summoner struct {
	Level         int64
	ProfileIconId int32
	Name          string
	LastRevision  time.Time
	TimeStamp     time.Time
}

type SummonerWithIds struct {
	Puuid         string
	AccountId     string
	SummonerId    string
	Level         int64
	ProfileIconId int32
	Name          string
	LastRevision  time.Time
	TimeStamp     time.Time
}
