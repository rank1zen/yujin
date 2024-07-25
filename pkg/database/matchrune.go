package database

type RuneSlot string

const (
	MainPath       RuneSlot = "main path"
	MainKeystone   RuneSlot = "main keystone"
	MainSlot1      RuneSlot = "main slot1"
	MainSlot2      RuneSlot = "main slot2"
	MainSlot3      RuneSlot = "main slot3"
	SecondaryPath  RuneSlot = "secondary path"
	SecondarySlot1 RuneSlot = "secondary slot1"
	SecondarySlot2 RuneSlot = "secondary slot2"
	ShardSlot1     RuneSlot = "shard slot1"
	ShardSlot2     RuneSlot = "shard slot2"
	ShardSlot3     RuneSlot = "shard slot3"
)

type MatchRuneFull struct {
	MainPath       int `db:"main_path"`
	MainKeystone   int `db:"main_keystone"`
	MainSlot1      int `db:"main_slot1"`
	MainSlot2      int `db:"main_slot2"`
	MainSlot3      int `db:"main_slot3"`
	SecondarySlot1 int `db:"secondary_slot1"`
	SecondarySlot2 int `db:"secondary_slot2"`
	ShardSlot1     int `db:"shard_slot1"`
	ShardSlot2     int `db:"shard_slot2"`
	ShardSlot3     int `db:"shard_slot3"`
}

func identifyRuneFull(m *MatchRuneFull, slot RuneSlot, runeID int) error {
	switch slot {
	case MainKeystone:
		m.MainKeystone = runeID
	case MainSlot1:
		m.MainSlot1 = runeID
	case MainSlot2:
		m.MainSlot2 = runeID
	case MainSlot3:
		m.MainSlot3 = runeID
	case SecondarySlot1:
		m.SecondarySlot1 = runeID
	case SecondarySlot2:
		m.SecondarySlot2 = runeID
	case ShardSlot1:
		m.ShardSlot1 = runeID
	case ShardSlot2:
		m.ShardSlot2 = runeID
	case ShardSlot3:
		m.ShardSlot3 = runeID
	}

	return nil
}

type MatchRuneSimple struct {
	MainKeystone  int
	SecondaryPath int
}

func identifyRuneSimple(m *MatchRuneSimple, slot RuneSlot, runeID int) error {
	switch slot {
	case MainKeystone:
		m.MainKeystone = runeID
	case SecondaryPath:
		m.SecondaryPath = runeID
	}

	return nil
}
