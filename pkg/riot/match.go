package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const endpointMatch = "/lol/match/v5"

type MatchDto struct {
	Metadata MatchMetadataDto `json:"metadata"`
	Info     MatchInfoDto     `json:"info"`
}

type MatchMetadataDto struct {
	DataVersion  string   `json:"dataVersion"`
	MatchId      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type MatchInfoDto struct {
	GameCreation int64 `json:"gameCreation"` // Unix timestamp for when the game is created on the game server (i.e., the loading screen).

	// Unix timestamp for when match starts on the game server.
	GameStartTimestamp int64 `json:"gameStartTimestamp"`
	// Unix timestamp for when match ends on the game server. This timestamp can occasionally
	// be significantly longer than when the match "ends". The most reliable way of determining
	// the timestamp for the end of the match would be to add the max time played of any
	// participant to the gameStartTimestamp. This field was added to match-v5 in patch 11.20 on Oct 5th, 2021.
	GameEndTimestamp int64 `json:"gameEndTimestamp"`

	// Prior to patch 11.20, this field returns the game length in milliseconds calculated
	// from gameEndTimestamp - gameStartTimestamp. Post patch 11.20, this field returns the max
	// timePlayed of any participant in the game in seconds, which makes the behavior of this
	// field consistent with that of match-v4. The best way to handling the change in this field
	// is to treat the value as milliseconds if the gameEndTimestamp field isn't in the response
	// and to treat the value as seconds if gameEndTimestamp is in the response.
	GameDuration int `json:"gameDuration"`

	GameVersion  string           `json:"gameVersion"`
	PlatformId   string           `json:"platformId"`
	Teams        []TeamDto        `json:"teams"`
	Participants []ParticipantDto `json:"participants"`
	// TournamentCode string `json:"tournamentCode"`
}

type TeamDto struct {
	Bans       []TeamBanDto  `json:"bans"`
	Objectives ObjectivesDto `json:"objectives"`
	TeamId     int           `json:"teamId"`
	Win        bool          `json:"win"`
}

type TeamBanDto struct {
	PickTurn   int `json:"pickTurn"`
	ChampionId int `json:"championId"`
}

type ObjectiveDto struct {
	First bool `json:"first"`
	Kills int  `json:"kills"`
}

type ObjectivesDto struct {
	Baron      ObjectiveDto `json:"baron"`
	Champion   ObjectiveDto `json:"champion"`
	Dragon     ObjectiveDto `json:"dragon"`
	Inhibitor  ObjectiveDto `json:"inhibitor"`
	RiftHerald ObjectiveDto `json:"riftHerald"`
	Tower      ObjectiveDto `json:"tower"`
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
	Description string       `json:"description"`
	Selections  []Selections `json:"selections"`
	Style       int          `json:"style"`
}

type ParticipantPerks struct {
	StatPerks *StatPerks `json:"statPerks"`
	Styles    []Styles   `json:"styles"`
}

type ParticipantDto struct {
	Assists         int `json:"assists"`
	BaronKills      int `json:"baronKills"`
	BountyLevel     int `json:"bountyLevel"`
	ChampExperience int `json:"champExperience"`
	ChampLevel      int `json:"champLevel"`
	// Prior to patch 11.4, on Feb 18th, 2021, this field returned invalid championIds.
	// We recommend determining the champion based on the championName field for matches played prior to patch 11.4.
	ChampionID   int    `json:"championId"`
	ChampionName string `json:"championName"`
	// This field is currently only utilized for Kayn's transformations.
	// (Legal values: 0 - None, 1 - Slayer, 2 - Assassin)
	ChampionTransform         int  `json:"championTransform"`
	ConsumablesPurchased      int  `json:"consumablesPurchased"`
	DamageDealtToBuildings    int  `json:"damageDealtToBuildings"`
	DamageDealtToObjectives   int  `json:"damageDealtToObjectives"`
	DamageDealtToTurrets      int  `json:"damageDealtToTurrets"`
	DamageSelfMitigated       int  `json:"damageSelfMitigated"`
	Deaths                    int  `json:"deaths"`
	DetectorWardsPlaced       int  `json:"detectorWardsPlaced"`
	DoubleKills               int  `json:"doubleKills"`
	DragonKills               int  `json:"dragonKills"`
	FirstBloodAssist          bool `json:"firstBloodAssist"`
	FirstBloodKill            bool `json:"firstBloodKill"`
	FirstTowerAssist          bool `json:"firstTowerAssist"`
	FirstTowerKill            bool `json:"firstTowerKill"`
	GameEndedInEarlySurrender bool `json:"gameEndedInEarlySurrender"`
	GameEndedInSurrender      bool `json:"gameEndedInSurrender"`
	GoldEarned                int  `json:"goldEarned"`
	GoldSpent                 int  `json:"goldSpent"`
	// Both individualPosition and teamPosition are computed by the game server and are
	// different versions of the most likely position played by a player. The individualPosition
	// is the best guess for which position the player actually played in isolation of
	// anything else. The teamPosition is the best guess for which position the player
	// actually played if we add the constraint that each team must have one top player, one
	// jungle, one middle, etc. Generally the recommendation is to use the teamPosition field
	// over the individualPosition field.
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
	// Both individualPosition and teamPosition are computed by the game server and are
	// different versions of the most likely position played by a player. The individualPosition
	// is the best guess for which position the player actually played in isolation of
	// anything else. The teamPosition is the best guess for which position the player
	// actually played if we add the constraint that each team must have one top player, one
	// jungle, one middle, etc. Generally the recommendation is to use the teamPosition field
	// over the individualPosition field.
	TeamPosition                   string `json:"teamPosition"`
	TimeCCingOthers                int    `json:"timeCCingOthers"`
	TimePlayed                     int    `json:"timePlayed"`
	TotalDamageDealt               int    `json:"totalDamageDealt"`
	TotalDamageDealtToChampions    int    `json:"totalDamageDealtToChampions"`
	TotalDamageShieldedOnTeammates int    `json:"totalDamageShieldedOnTeammates"`
	TotalDamageTaken               int    `json:"totalDamageTaken"`
	TotalHeal                      int    `json:"totalHeal"`
	TotalHealsOnTeammates          int    `json:"totalHealsOnTeammates"`
	TotalMinionsKilled             int    `json:"totalMinionsKilled"`
	TotalTimeCCDealt               int    `json:"totalTimeCCDealt"`
	TotalTimeSpentDead             int    `json:"totalTimeSpentDead"`
	TotalUnitsHealed               int    `json:"totalUnitsHealed"`
	TripleKills                    int    `json:"tripleKills"`
	TrueDamageDealt                int    `json:"trueDamageDealt"`
	TrueDamageDealtToChampions     int    `json:"trueDamageDealtToChampions"`
	TrueDamageTaken                int    `json:"trueDamageTaken"`
	TurretKills                    int    `json:"turretKills"`
	TurretTakedowns                int    `json:"turretTakedowns"`
	TurretsLost                    int    `json:"turretsLost"`
	UnrealKills                    int    `json:"unrealKills"`
	VisionScore                    int    `json:"visionScore"`
	VisionWardsBoughtInGame        int    `json:"visionWardsBoughtInGame"`
	WardsKilled                    int    `json:"wardsKilled"`
	WardsPlaced                    int    `json:"wardsPlaced"`
	Win                            bool   `json:"win"`
}

func listByPuuid(ctx context.Context, doer Doer, puuid string, start int, count int) ([]string, error) {
	url := fmt.Sprintf(endpointMatch+"/matches/by-puuid/%s/ids?queue=420&start=%d&count=%d", puuid, start, count)
	req, err := NewRiotRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := doer.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("riot api: %s", res.Status)
	}

	var ids []string
	err = json.NewDecoder(res.Body).Decode(&ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

// Get a match by match id
func matchById(ctx context.Context, doer Doer, matchID string) (*MatchDto, error) {
	url := fmt.Sprintf(endpointMatch+"/matches/%s", matchID)
	req, err := NewRiotRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := doer.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("riot api: %s", res.Status)
	}

	var match MatchDto
	err = json.NewDecoder(res.Body).Decode(&match)
	if err != nil {
		return nil, err
	}

	return &match, nil
}
