package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchV5Query struct {
	db *pgxpool.Pool
}

func NewMatchV5Query(pool *pgxpool.Pool) *MatchV5Query {
	return &MatchV5Query{
		db: pool,
	}
}

func (q *MatchV5Query) InsertMatch(ctx context.Context, arg *MatchRecordArg) (string, error) {
	_ = `
	INSERT INTO match_v5
	(match_id, runes)
	VALUES ($1, 
		ROW(ROW($2, $3, $4))
	)
	RETURNING match_id
	`

	return "", nil
}

func (q *MatchV5Query) InsertMatchTeam(ctx context.Context, arg *MatchTeamRecordArg) (string, error) {
	var id string
	err := q.db.QueryRow(ctx, `
		INSERT INTO match_team_record
		(record_date, match_id, team_id, objectives, bans)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING record_id
	`, arg.RecordDate, arg.MatchId, arg.TeamId, arg.Objective, arg.Bans).Scan(&id)
	if err != nil {
		return "", err
	}
	
	return id, nil
}

func (q *MatchV5Query) InsertMatchPlayer(ctx context.Context, arg *MatchParticipantRecordArg) (string, error) {
	var uuid string
	err := q.db.QueryRow(ctx, `
		INSERT INTO match_participant
		(win, position, kills, deaths, assists, creep_score, gold_earned, champion_id, champion_name, champion_level,
		item0, item1, item2, item3, item4, item5, item6, vision_score, wards_placed, control_wards_placed,
		first_blood_assist, first_tower_assist, turret_takedowns, physical_damage_dealt_to_champions,
		magic_damage_dealt_to_champions, true_damage_dealt_to_champions, total_damage_dealt_to_champions,
		total_damage_taken, total_heals_on_teammates)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`,
		arg.Win,
		arg.Position,
		arg.Kills,
		arg.Deaths,
		arg.Assists,
		arg.CS,
		arg.Gold,
		arg.ChampId,
		arg.ChampName,
		arg.ChampLevel,
		arg.Item0,
		arg.Item1,
		arg.Item2,
		arg.Item3,
		arg.Item4,
		arg.Item5,
		arg.Item6,
		arg.VisionScore,
		arg.WardsPlaced,
		arg.ControlWardsPlaced,
		arg.FirstBloodAssist,
		arg.FirstTowerAssist,
		arg.TurretTakeDowns,
		arg.PhysicalDamageDealtToChampions,
		arg.MagicDamageDealtToChampions,
		arg.TrueDamageDealtToChampions,
		arg.TotalDamageDealtToChampions,
		arg.TotalDamageTaken,
		arg.TotalHealsOnTeammates,
	).Scan(&uuid)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (q *MatchV5Query) SelectMatch(ctx context.Context, id string) (MatchRecordArg, error) {
	rows, _ := q.db.Query(ctx, `
		SELECT * FROM match
		WHERE match_id = $1
	`, id)
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[MatchRecordArg])
}

func (q *MatchV5Query) SelectMatchPlayer(ctx context.Context, id string) {
}

func (q *MatchV5Query) SelectMatchTeam(ctx context.Context, id string) ([]*MatchTeamRecord, error) {
	rows, _ := q.db.Query(ctx, `
		SELECT * FROM match_team_record
		WHERE match_id = $1
	`, id)
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[MatchTeamRecord])
}
