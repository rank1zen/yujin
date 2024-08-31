package database

import (
	"context"
	"time"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/internal/riot"
)

type MatchKillEvent struct {
	Timestamp time.Duration `db:"timestamp"`
	PosX      int           `db:"pos_x"`
	PosY      int           `db:"pos_y"`
	Bounty    int           `db:"bounty"`
	Shutdown  int           `db:"shutdown_bounty"`
	KillerID  int           `db:"killer_id"`
	VictimID  int           `db:"victim_id"`
}

func matchKillEvent(batch *pgx.Batch, m *riot.MatchTimeline) {
	panic("broken")
	for _, frame := range m.Info.Frames {
		for _, event := range frame.Events {
			switch event.EventType {
			case "HI":
				batch.Queue(`
				insert into match_champion_kill_event_records
				()
				VALUES
				()
				`)
			}
		}
	}
}

type MatchKillEventList []*MatchKillEvent

type MatchItemEvent struct {
	Timestamp     time.Duration `db:"timestamp"`
	ParticipantID int
	SkillSlot     int
}

type MatchItemEventList []*MatchItemEvent

type MatchSkillEvent struct {
	Timestamp     time.Duration `db:"timestamp"`
	ParticipantID int
	ItemID        int
}

type MatchSkillEventList []*MatchSkillEvent

type SummonerMatchItemEvent struct {
	MatchId       RiotMatchId
	Timestamp     time.Duration
	ParticipantId int
	ItemId        int
}

func (m SummonerMatchItemEvent) GetItemIconUrl() {}

type SummonerMatchKillEvent struct{}

type SummonerMatchSkillEvent struct {
	MatchId       RiotMatchId   `db:"match_id"`
	Timestamp     time.Duration `db:"timestamp"`
	ParticipantId int           `db:"participant_id"`
	SkillSlot     int           `db:"skill_slot"`
}

func (m SummonerMatchSkillEvent) GetItemIconUrl() templ.SafeURL {
	// FIXME: Please
	return templ.URL("a")
}

func (m SummonerMatchSkillEvent) GetTimestamp() string {
	return m.Timestamp.String()
}

type SummonerMatchSkillEventList []*SummonerMatchSkillEvent

func (db *DB) get(ctx context.Context, matchID string, puuid string) ([]MatchItemEvent, error) {
	querySelect(ctx, db.pool, `
	SELECT * FROM match_spell_event_records
	WHERE match_id = $1 AND participant_id = $2,
	`, []any{matchID, puuid}, pgx.RowToStructByNameLax[string])

	return nil, nil
}
