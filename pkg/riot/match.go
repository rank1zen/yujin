package riot

import (
	"context"
	"fmt"
)

type Match struct {
	Info     *MatchInfo     `json:"info"`
	Metadata *MatchMetadata `json:"metadata"`
}

type MatchMetadata struct {
	DataVersion  string   `json:"dataVersion"`
	MatchId      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type MatchInfo struct {
	GameCreation       int64          `json:"gameCreation"`
	GameStartTimestamp int64          `json:"gameStartTimestamp"`
	GameEndTimestamp   int64          `json:"gameEndTimestamp"`
	GameDuration       int            `json:"gameDuration"`
	GameVersion        string         `json:"gameVersion"`
	PlatformId         string         `json:"platformId"`
	Teams              []*Team        `json:"teams"`
	Participants       []*Participant `json:"participants"`
}

type Team struct {
	Bans       []*TeamBan  `json:"bans"`
	Objectives *Objectives `json:"objectives"`
	TeamId     int         `json:"teamId"`
	Win        bool        `json:"win"`
}

type TeamBan struct {
	PickTurn   int `json:"pickTurn"`
	ChampionId int `json:"championId"`
}

type Objective struct {
	First bool `json:"first"`
	Kills int  `json:"kills"`
}

type Objectives struct {
	Baron      *Objective `json:"baron"`
	Champion   *Objective `json:"champion"`
	Dragon     *Objective `json:"dragon"`
	Inhibitor  *Objective `json:"inhibitor"`
	RiftHerald *Objective `json:"riftHerald"`
	Tower      *Objective `json:"tower"`
}

type StatPerks struct {
	Defense int `json:"defense"`
	Flex    int `json:"flex"`
	Offense int `json:"offense"`
}

type Selections struct {
	Perk int `json:"perk"`
	Var1 int `json:"var1"`
	Var2 int `json:"var2"`
	Var3 int `json:"var3"`
}

type Styles struct {
	Description string        `json:"description"`
	Selections  []*Selections `json:"selections"`
	Style       int           `json:"style"`
}

type ParticipantPerks struct {
	StatPerks *StatPerks `json:"statPerks"`
	Styles    []*Styles  `json:"styles"`
}

type Participant struct {
	Assists                        int               `json:"assists"`
	BaronKills                     int               `json:"baronKills"`
	BountyLevel                    int               `json:"bountyLevel"`
	ChampExperience                int               `json:"champExperience"`
	ChampLevel                     int               `json:"champLevel"`
	ChampionID                     int               `json:"championId"`
	ChampionName                   string            `json:"championName"`
	ChampionTransform              int               `json:"championTransform"`
	ConsumablesPurchased           int               `json:"consumablesPurchased"`
	DamageDealtToBuildings         int               `json:"damageDealtToBuildings"`
	DamageDealtToObjectives        int               `json:"damageDealtToObjectives"`
	DamageDealtToTurrets           int               `json:"damageDealtToTurrets"`
	DamageSelfMitigated            int               `json:"damageSelfMitigated"`
	Deaths                         int               `json:"deaths"`
	DetectorWardsPlaced            int               `json:"detectorWardsPlaced"`
	DoubleKills                    int               `json:"doubleKills"`
	DragonKills                    int               `json:"dragonKills"`
	FirstBloodAssist               bool              `json:"firstBloodAssist"`
	FirstBloodKill                 bool              `json:"firstBloodKill"`
	FirstTowerAssist               bool              `json:"firstTowerAssist"`
	FirstTowerKill                 bool              `json:"firstTowerKill"`
	GameEndedInEarlySurrender      bool              `json:"gameEndedInEarlySurrender"`
	GameEndedInSurrender           bool              `json:"gameEndedInSurrender"`
	GoldEarned                     int               `json:"goldEarned"`
	GoldSpent                      int               `json:"goldSpent"`
	IndividualPosition             string            `json:"individualPosition"`
	InhibitorKills                 int               `json:"inhibitorKills"`
	InhibitorTakedowns             int               `json:"inhibitorTakedowns"`
	InhibitorsLost                 int               `json:"inhibitorsLost"`
	Item0                          int               `json:"item0"`
	Item1                          int               `json:"item1"`
	Item2                          int               `json:"item2"`
	Item3                          int               `json:"item3"`
	Item4                          int               `json:"item4"`
	Item5                          int               `json:"item5"`
	Item6                          int               `json:"item6"`
	ItemsPurchased                 int               `json:"itemsPurchased"`
	KillingSprees                  int               `json:"killingSprees"`
	Kills                          int               `json:"kills"`
	Lane                           string            `json:"lane"`
	LargestCriticalStrike          int               `json:"largestCriticalStrike"`
	LargestKillingSpree            int               `json:"largestKillingSpree"`
	LargestMultiKill               int               `json:"largestMultiKill"`
	LongestTimeSpentLiving         int               `json:"longestTimeSpentLiving"`
	MagicDamageDealt               int               `json:"magicDamageDealt"`
	MagicDamageDealtToChampions    int               `json:"magicDamageDealtToChampions"`
	MagicDamageTaken               int               `json:"magicDamageTaken"`
	NeutralMinionsKilled           int               `json:"neutralMinionsKilled"`
	NexusKills                     int               `json:"nexusKills"`
	NexusLost                      int               `json:"nexusLost"`
	NexusTakedowns                 int               `json:"nexusTakedowns"`
	ObjectivesStolen               int               `json:"objectivesStolen"`
	ObjectivesStolenAssists        int               `json:"objectivesStolenAssists"`
	ParticipantID                  int               `json:"participantId"`
	PentaKills                     int               `json:"pentaKills"`
	Perks                          *ParticipantPerks `json:"perks"`
	PhysicalDamageDealt            int               `json:"physicalDamageDealt"`
	PhysicalDamageDealtToChampions int               `json:"physicalDamageDealtToChampions"`
	PhysicalDamageTaken            int               `json:"physicalDamageTaken"`
	ProfileIcon                    int               `json:"profileIcon"`
	PUUID                          string            `json:"puuid"`
	QuadraKills                    int               `json:"quadraKills"`
	RiotIDGameName                 string            `json:"riotIdGameName"`
	RiotIDName                     string            `json:"riotIdName"`
	RiotIDTagline                  string            `json:"riotIdTagline"`
	Role                           string            `json:"role"`
	SightWardsBoughtInGame         int               `json:"sightWardsBoughtInGame"`
	Spell1Casts                    int               `json:"spell1Casts"`
	Spell2Casts                    int               `json:"spell2Casts"`
	Spell3Casts                    int               `json:"spell3Casts"`
	Spell4Casts                    int               `json:"spell4Casts"`
	Summoner1Casts                 int               `json:"summoner1Casts"`
	Summoner1ID                    int               `json:"summoner1Id"`
	Summoner2Casts                 int               `json:"summoner2Casts"`
	Summoner2ID                    int               `json:"summoner2Id"`
	SummonerID                     string            `json:"summonerId"`
	SummonerLevel                  int               `json:"summonerLevel"`
	SummonerName                   string            `json:"summonerName"`
	TeamEarlySurrendered           bool              `json:"teamEarlySurrendered"`
	TeamID                         int               `json:"teamId"`
	TeamPosition                   string            `json:"teamPosition"`
	TimeCCingOthers                int               `json:"timeCCingOthers"`
	TimePlayed                     int               `json:"timePlayed"`
	TotalDamageDealt               int               `json:"totalDamageDealt"`
	TotalDamageDealtToChampions    int               `json:"totalDamageDealtToChampions"`
	TotalDamageShieldedOnTeammates int               `json:"totalDamageShieldedOnTeammates"`
	TotalDamageTaken               int               `json:"totalDamageTaken"`
	TotalHeal                      int               `json:"totalHeal"`
	TotalHealsOnTeammates          int               `json:"totalHealsOnTeammates"`
	TotalMinionsKilled             int               `json:"totalMinionsKilled"`
	TotalTimeCCDealt               int               `json:"totalTimeCCDealt"`
	TotalTimeSpentDead             int               `json:"totalTimeSpentDead"`
	TotalUnitsHealed               int               `json:"totalUnitsHealed"`
	TripleKills                    int               `json:"tripleKills"`
	TrueDamageDealt                int               `json:"trueDamageDealt"`
	TrueDamageDealtToChampions     int               `json:"trueDamageDealtToChampions"`
	TrueDamageTaken                int               `json:"trueDamageTaken"`
	TurretKills                    int               `json:"turretKills"`
	TurretTakedowns                int               `json:"turretTakedowns"`
	TurretsLost                    int               `json:"turretsLost"`
	UnrealKills                    int               `json:"unrealKills"`
	VisionScore                    int               `json:"visionScore"`
	VisionWardsBoughtInGame        int               `json:"visionWardsBoughtInGame"`
	WardsKilled                    int               `json:"wardsKilled"`
	WardsPlaced                    int               `json:"wardsPlaced"`
	Win                            bool              `json:"win"`
}

// TODO: implement
type MatchTimeline struct{}

// Get a list of match ids by puuid (ONLY RANKED SOLOQ, queueId 420)
func (s *Client) GetMatchHistory(ctx context.Context, puuid string, start int, count int) ([]string, error) {
	u := fmt.Sprintf(defaultBaseURLTemplate+"/lol/match/v5/matches/by-puuid/%s/ids?queue=420&start=%d&count=%d", "americas", puuid, start, count)
	req := NewRequest(
		WithToken2(),
		WithURL(u),
	)

	var ids []string
	err := s.Do(ctx, req, &ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

// Get a match by match id
func (c *Client) GetMatch(ctx context.Context, matchID string) (*Match, error) {
	u := fmt.Sprintf(defaultBaseURLTemplate+"/lol/match/v5/matches/%s", "americas", matchID)
	req := NewRequest(
		WithURL(u),
		WithToken2(),
	)

	var match Match
	err := c.Do(ctx, req, &match)
	if err != nil {
		return nil, err
	}

	return &match, nil
}

// Get a match timeline by match id
// TODO: implement
func (c *Client) GetMatchTimeline(ctx context.Context, matchID string) (*MatchTimeline, error) {
	return nil, nil
}
