package main

import "time"

type SummonerRecordBody struct {
	RecordDate    time.Time `json:"record_date" validate:"required"`
	AccountId     string    `json:"account_id"`
	ProfileIconId int32     `json:"profile_icon_id"`
	RevisionDate  int64     `json:"revision_date"`
	Name          string    `json:"name"`
	Puuid         string    `json:"puuid"`
	SummonerLevel int64     `json:"summoner_level"`
}

type MatchBody struct {
	MatchId string `json:"match_id" validate:"required"`
}

type SummonerProfileQuery struct {
}

type SummonerProfileBody struct {
	Name       string `json:"name"`
	Puuid      string `json:"puuid"`
	AccountId  string `json:"account_id"`
	SummonerId string `json:"summoner_id"`
}
