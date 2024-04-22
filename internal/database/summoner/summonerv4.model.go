package summoner

import "time"

type SummonerRecordDatum struct {
	RecordId   string    `db:"record_id"`
	RecordDate time.Time `db:"record_date"`

	Puuid         string `db:"puuid"`
	AccountId     string `db:"account_id"`
	SummonerId    string `db:"id"`
	Name          string `db:"name"`
	ProfileIconId int32  `db:"profile_icon_id"`
	SummonerLevel int32  `db:"summoner_level"`
	RevisionDate  int64  `db:"revision_date"`
}

type SummonerRecordArg struct {
	RecordDate time.Time

	Puuid         string
	AccountId     string
	SummonerId    string
	Name          string
	ProfileIconId int32
	SummonerLevel int32
	RevisionDate  int64
}

type SummonerRecordFilter struct {
	Field   string
	Value   string
	DateMin time.Time `json:"timestamp_gt,omitempty"`
	DateMax time.Time `json:"timestamp_lt,omitempty"`

	SortOrder string
	Offset    uint64
	Limit     uint64
}

type SummonerRecordCountFilter struct {
	Field   string
	Value   string
	DateMin time.Time `json:"timestamp_gte,omitempty"`
	DateMax time.Time `json:"timestamp_lte,omitempty"`
}
