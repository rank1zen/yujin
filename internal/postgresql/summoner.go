package postgresql

import (
	"context"

	"github.com/rank1zen/yujin/internal"
	"gorm.io/gorm"
)

type SummonerDA struct {
	db *gorm.DB
}

func NewSummonerDA(db *gorm.DB) *SummonerDA {
	return &SummonerDA{
		db: db,
	}
}

func (s *SummonerDA) Create(ctx context.Context, summoner Summoner) error {
	tx := s.db.Create(summoner)
	if tx.Error != nil {
		return internal.WrapErrorf(tx.Error, internal.ErrorCodeUnknown, "insert summoner")
	}
	return nil
}

func (s *SummonerDA) Find(ctx context.Context, puuid string) ([]Summoner, error) {
	var sum []Summoner
	tx := s.db.Where("puuid = ?", puuid).Find(&sum)

	if tx.Error != nil {
		return nil, internal.WrapErrorf(tx.Error, internal.ErrorCodeUnknown, "summoner not found")
	}

	return sum, nil
}

func (s *SummonerDA) Newest(ctx context.Context, puuid string) (Summoner, error) {
	var sum Summoner
	tx := s.db.Where("puuid = ?", puuid).First(&sum)
	if tx.Error != nil {
		return Summoner{}, internal.WrapErrorf(tx.Error, internal.ErrorCodeUnknown, "summoner not found")
	}
	return sum, nil
}

type RiotRegion string

const (
	RegionNA  RiotRegion = "na"
	RegionEUW RiotRegion = "euw"
	RegionKR  RiotRegion = "kr"
)

type Summoner struct {
	gorm.Model
	Region        RiotRegion
	Puuid         string
	AccountId     string
	SummonerId    string
	Level         int64
	ProfileIconId int32
	Name          string
	LastRevision  int64
	TimeStamp     int64
}
