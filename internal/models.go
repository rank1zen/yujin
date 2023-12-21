package internal

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rank1zen/yujin/internal/postgresql/db"
)

type Summoner struct {
	Level         int64
	ProfileIconId int32
	Name          string
	LastRevision  time.Time
	TimeStamp     time.Time
}

type SummonerWithIds struct {
	Puuid         string
	AccountId     string
	SummonerId    string
	Level         int64
	ProfileIconId int32
	Name          string
	LastRevision  time.Time
	TimeStamp     time.Time
}

func newTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

func ParseUUID(s string) (pgtype.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}

	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}, nil
}

func UUIDString(uuid pgtype.UUID) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid.Bytes[0:4], uuid.Bytes[4:6], uuid.Bytes[6:8], uuid.Bytes[8:10], uuid.Bytes[10:16])
}

func (s *SummonerWithIds) CastToDB() db.InsertSummonerParams {
	return db.InsertSummonerParams{
		Puuid:         s.Puuid,
		AccountID:     s.AccountId,
		SummonerID:    s.SummonerId,
		Level:         s.Level,
		ProfileIconID: s.ProfileIconId,
		Name:          s.Name,
		LastRevision:  newTimestamp(s.LastRevision),
		TimeStamp:     newTimestamp(s.TimeStamp),
	}
}
