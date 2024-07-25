package database

import (
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/pkg/riot"
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

type MatchSkillEvent struct {
	Timestamp     time.Duration `db:"timestamp"`
	ParticipantID int
	ItemID        int
}
