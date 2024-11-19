package riot

type PUUID string

func (m PUUID) String() string {
	return string(m)
}

type SummonerID string

type MatchID string

func (m MatchID) String() string {
	return string(m)
}

type TeamID int32

type ParticipantID int32
