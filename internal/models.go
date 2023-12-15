package internal

import "time"

type Summoner struct {
	Puuid         string
	AccountId     string
	SummonerId    string
	Level         int
	ProfileIconId int
	Name          string
	LastRevision  time.Time
	TimeStamp     time.Time
}
