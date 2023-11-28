package postgresql

import (
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/postgresql/db"
	"gorm.io/gorm"
)

// Summoner CRUD: Create, Find, FindRecent

// Repository injection
type SummonerDA struct {
	db *gorm.DB 
}

func (s *SummonerDA) Create(summoner *internal.MSummoner) error {
	s.db.Create(&summoner)
	// Something like this
	return nil
}

func (s *SummonerDA) Find(params *internal.CreateSummonerParams) (internal.MSummoner, error) {
	return nil, nil;
}

func (s *SummonerDA) FindRecent(puuid string) (internal.MSummoner, error) {
	r := s.db.First(&db.DBSummoner, "a")
	return r, nil
}
