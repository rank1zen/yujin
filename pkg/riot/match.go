package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Get a list of match ids by puuid (ONLY RANKED SOLOQ, queueId 420)
func (s *Client) GetMatchHistory(ctx context.Context, puuid string, start int, count int) ([]string, error) {
	u := fmt.Sprintf(defaultAmerBaseURL+"/lol/match/v5/matches/by-puuid/%s/ids?queue=420&start=%d&count=%d", puuid, start, count)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("X-Riot-Token", os.Getenv("RIOT_API_KEY"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var ids []string
	err = json.NewDecoder(res.Body).Decode(&ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

type Match struct {
	Info     *MatchInfo     `json:"info"`
	Metadata *MatchMetadata `json:"metadata"`
}

type MatchList []*Match

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

// Get a match by match id
func (c *Client) GetMatch(ctx context.Context, matchID string) (*Match, error) {
	u := fmt.Sprintf(defaultAmerBaseURL+"/lol/match/v5/matches/%s", matchID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("X-Riot-Token", os.Getenv("RIOT_API_KEY"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var m Match
	err = json.NewDecoder(res.Body).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("decoding: %w", err)
	}

	return &m, nil
}

// TODO: these are todos

type MatchTimeline struct {
	Metadata *MetadataTimeLine
	Info *InfoTimeLine
}

type MetadataTimeLine struct {}

type InfoTimeLine struct {
	EndOfGameResult 	string
	FrameInterval 	int64 	
	GameId 	int64 	
	Participants 	[]*ParticipantTimeLine
	Frames 	[]*FramesTimeLine
}

type ParticipantTimeLine struct { }

type FramesTimeLine struct {
	ParticipantFrames map[string]*ParticipantFrame `json:"participantFrames"`
	Events            []*MatchEvent                `json:"events"`
	Timestamp         int                          `json:"timestamp"`
}

type MatchPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// ParticipantFrame contains information about a participant in a game at a single timestamp
type ParticipantFrame struct {
	Position            *MatchPosition `json:"position"`
	TotalGold           int            `json:"totalGold"`
	TeamScore           int            `json:"teamScore"`
	ParticipantID       int            `json:"participantId"`
	Level               int            `json:"level"`
	CurrentGold         int            `json:"currentGold"`
	MinionsKilled       int            `json:"minionsKilled"`
	DominionScore       int            `json:"dominionScore"`
	XP                  int            `json:"xp"`
	JungleMinionsKilled int            `json:"jungleMinionsKilled"`
}

// MatchEventType is the type of an event
type MatchEventType string

// All legal value for match event types
const (
	MatchEventTypeChampionKill     MatchEventType = "CHAMPION_KILL"
	MatchEventTypeWardPlaced       MatchEventType = "WARD_PLACED"
	MatchEventTypeWardKill         MatchEventType = "WARD_KILL"
	MatchEventTypeBuildingKill     MatchEventType = "BUILDING_KILL"
	MatchEventTypeEliteMonsterKill MatchEventType = "ELITE_MONSTER_KILL"
	MatchEventTypeItemPurchased    MatchEventType = "ITEM_PURCHASED"
	MatchEventTypeItemSold         MatchEventType = "ITEM_SOLD"
	MatchEventTypeItemDestroyed    MatchEventType = "ITEM_DESTROYED"
	MatchEventTypeItemUndo         MatchEventType = "ITEM_UNDO"
	MatchEventTypeSkillLevelUp     MatchEventType = "SKILL_LEVEL_UP"
	MatchEventTypeAscendedEvent    MatchEventType = "ASCENDED_EVENT"
	MatchEventTypeCapturePoint     MatchEventType = "CAPTURE_POINT"
	MatchEventTypePoroKingSummon   MatchEventType = "PORO_KING_SUMMON"
)

var (
	// MatchEventTypes is a list of all available match events
	MatchEventTypes = []MatchEventType{
		MatchEventTypeChampionKill,
		MatchEventTypeWardPlaced,
		MatchEventTypeWardKill,
		MatchEventTypeBuildingKill,
		MatchEventTypeEliteMonsterKill,
		MatchEventTypeItemPurchased,
		MatchEventTypeItemSold,
		MatchEventTypeItemDestroyed,
		MatchEventTypeItemUndo,
		MatchEventTypeSkillLevelUp,
		MatchEventTypeAscendedEvent,
		MatchEventTypeCapturePoint,
		MatchEventTypePoroKingSummon,
	}
)

// MatchEvent is an event in a match at a certain timestamp
type MatchEvent struct {
	Type                    *MatchEventType `json:"type"`
	Position                *MatchPosition  `json:"position"`
	LevelUpType             string          `json:"levelUpType"`
	AscendedType            string          `json:"ascendedType"`
	TowerType               string          `json:"towerType"`
	EventType               string          `json:"eventType"`
	PointCaptured           string          `json:"pointCaptured"`
	WardType                string          `json:"wardType"`
	MonsterType             string          `json:"monsterType"`
	BuildingType            string          `json:"buildingType"`
	LaneType                string          `json:"laneType"`
	MonsterSubType          string          `json:"monsterSubType"`
	AssistingParticipantIDs []int           `json:"assistingParticipantIds"`
	Timestamp               int             `json:"timestamp"`
	AfterID                 int             `json:"afterId"`
	VictimID                int             `json:"victimId"`
	SkillSlot               int             `json:"skillSlot"`
	ItemID                  int             `json:"itemId"`
	ParticipantID           int             `json:"participantId"`
	TeamID                  int             `json:"teamId"`
	CreatorID               int             `json:"creatorId"`
	KillerID                int             `json:"killerId"`
	BeforeID                int             `json:"beforeId"`
}

// Get a match timeline by match id
func (c *Client) GetMatchTimeline(ctx context.Context, matchID string) (*MatchTimeline, error) {
	u := fmt.Sprintf(defaultAmerBaseURL+"/lol/match/v5/matches/%s/timeline", matchID)
	req := NewRequest(
		WithURL(u),
		WithToken2(),
	)

	var timeline MatchTimeline
	err := c.Do(ctx, req, &timeline)
	if err != nil {
		return nil, err
	}

	return &timeline, nil
}
