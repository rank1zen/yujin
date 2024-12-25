package internal

import "time"

type LiveMatchParticipant struct {
	Puuid          PUUID
	TeamID         TeamID
	SummonerID     SummonerID
	Champion       ChampionID
	Runes          Runes
	BannedChampion *ChampionID
	Summoners      SummonersIDs
	// ProfileIconID ProfileIconID
}

type LiveMatchParticipantList [10]LiveMatchParticipant

type LiveMatchTeamList [5]LiveMatchParticipant

type LiveMatch struct {
	GameStartTime time.Time
	GameLength    time.Duration
	Participant   LiveMatchParticipantList
}

func (m LiveMatch) GetRed() LiveMatchParticipant {
	return m.Participant[0]
}
