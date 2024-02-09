package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchV5Query struct {
	db *pgxpool.Pool
}

func newMatchV5Query(pool *pgxpool.Pool) *MatchV5Query {
	return &MatchV5Query{
		db: pool,
	}
}

type Match struct {
	Id        string
	Runes Perks
	StartTs   int
	Duration  int64
	Surrender bool
	Patch     string
	Puuids    []string

	//RedSide  Team
	//BlueSide Team
	Players []Participant
}

type Team struct {
	Bans       []Ban
	Baron      Objective
	Dragon     Objective
	RiftHerald Objective
	Champion   Objective
	Tower      Objective
	Inhibitor  Objective
}

type Ban struct {
	ChampionId int
	PickTurn   int
}

type Objective struct {
	First bool
	Kills int
}

type Participant struct {
	Id            int
	Puuid         string
	SummonerLevel int
	SummonerName  string
	Win           bool
	Position      string

	Kills      int
	Deaths     int
	Assists    int
	CS         int
	Gold       int
	ChampId    int
	ChampName  int
	ChampLevel int
	Item0      int
	Item1      int
	Item2      int
	Item3      int
	Item4      int
	Item5      int
	Item6      int

	VisionScore                    int
	WardsPlaced                    int
	ControlWardsPlaced             int
	FirstBloodAssist               bool
	FirstTowerAssist               bool
	TurretTakeDowns                int
	PhysicalDamageDealtToChampions int
	MagicDamageDealtToChampions    int
	TrueDamageDealtToChampions     int
	TotalDamageDealtToChampions    int
	TotalDamageTaken               int
	TotalHealsOnTeammates          int

	Runes Perks
}

type Perks struct {
	Styles    []PerkStyle
	StatPerks PerkStats
}

type PerkStyle struct {
	Desc       string
	Style      int
	Selections []PerkStyleSelection
}

type PerkStyleSelection struct {
	Perk int
	Var1 int
	Var2 int
	Var3 int
}

type PerkStats struct {
	Defense int32
	Flex    int32
	Offense int32
}

func (q *MatchV5Query) SelectMatch(ctx context.Context, id string) (Match, error) {
	query := `
	SELECT * FROM match_v5
	WHERE match_id = $1
	`

	rows, _ := q.db.Query(ctx, query, id)
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Match])
}

func (q *MatchV5Query) InsertMatch(ctx context.Context, arg *Match) (string, error) {
	query := `
	INSERT INTO match_v5
	(match_id, runes)
	VALUES ($1, 
		ROW(ROW($2, $3, $4))
	)
	RETURNING match_id
	`

	var id string
	err := q.db.QueryRow(ctx, query, arg.Id, arg.Runes.StatPerks.Defense, arg.Runes.StatPerks.Flex, arg.Runes.StatPerks.Offense).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
