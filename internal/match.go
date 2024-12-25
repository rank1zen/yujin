package internal

import "time"

// Match is strictly ranked, soloq, 5v5, summoners rift.
type Match struct {
	ID              MatchID
	Patch           GameVersion
	CreateTimestamp time.Time
	StartTimestamp  time.Time
	EndTimestamp    time.Time
	Duration        time.Duration
	EndOfGameResult string
	Participants    MatchParticipants
}

func (m Match) GetParticipants() [10]Participant {
	return m.Participants
}

// GetBlueRed returns Blue side Red side
func (m Match) GetBlueRed() ([5]Participant, [5]Participant) {
	return [5]Participant(m.Participants[:5]), [5]Participant(m.Participants[:5])
}

type ParticipantID int

type Participant struct {
	ID       ParticipantID
	Puuid    PUUID
	Match    MatchID
	Team     TeamID
	Summoner SummonerID

	SummonerLevel             int
	SummonerName              string
	RiotIDGameName            string
	RiotIDName                string
	RiotIDTagline             string
	ChampionLevel             int
	ChampionID                ChampionID
	ChampionName              string
	GameEndedInEarlySurrender bool
	GameEndedInSurrender      bool
	Items                     ItemIDs
	Runes                     Runes
	Role                      string
	Summoners                 SummonersIDs
	TeamEarlySurrendered      bool
	TeamPosition              string
	TimePlayed                int
	Win                       bool

	// Post game stats

	Assists                        int
	DamageDealtToBuildings         int
	DamageDealtToObjectives        int
	DamageDealtToTurrets           int
	DamageSelfMitigated            int
	Deaths                         int
	DetectorWardsPlaced            int
	FirstBloodAssist               bool
	FirstBloodKill                 bool
	FirstTowerAssist               bool
	FirstTowerKill                 bool
	GoldEarned                     int
	GoldSpent                      int
	IndividualPosition             string
	InhibitorKills                 int
	InhibitorTakedowns             int
	InhibitorsLost                 int
	Kills                          int
	MagicDamageDealt               int
	MagicDamageDealtToChampions    int
	MagicDamageTaken               int
	PhysicalDamageDealt            int
	PhysicalDamageDealtToChampions int
	PhysicalDamageTaken            int
	SightWardsBoughtInGame         int
	TotalDamageDealt               int
	TotalDamageDealtToChampions    int
	TotalDamageShieldedOnTeammates int
	TotalDamageTaken               int
	TotalHeal                      int
	TotalHealsOnTeammates          int
	TotalMinionsKilled             int
	TrueDamageDealt                int
	TrueDamageDealtToChampions     int
	TrueDamageTaken                int
	VisionScore                    int
	VisionWardsBoughtInGame        int
	WardsKilled                    int
	WardsPlaced                    int
	NeutralMinionsKilled           int
}

type MatchParticipants [10]Participant

type TeamID int

type Team struct {
	ID      TeamID
	MatchID MatchID
	Win     bool
}

type MatchTeams [2]Team
