package postgresql

import (
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/postgresql/db"
	"gorm.io/gorm"
)

// Summoner CRUD: Create, Find, FindRecent

// Repository injection
type SummonerDA struct {
	q *gorm.DB 
}

func NewSummonerDA(q *gorm.DB) *SummonerDA {
	return &SummonerDA{
		q: q,
	}
}

func (s *SummonerDA) Create(summoner *internal.Summoner) error {
	ctx := s.q.Create(&summoner)
	if ctx.Error != nil {
		return internal.WrapErrorf(ctx.Error, internal.ErrorCodeUnknown, "create summoner")
	}
	return nil
}

func (s *SummonerDA) Find(params *internal.FindSummonerParams) ([]internal.Summoner, error) {
	var summoners []db.Summoner
	s.q.Limit(params.Limit).Offset(params.Limit).Find(&summoners)

	return nil, nil;
}

func (s *SummonerDA) FindRecent(puuid string) (internal.Summoner, error) {
	var summoner db.Summoner
	ctx := s.q.Where("puu_id = ?", puuid).First(&summoner)
	
	if ctx.Error != nil {
		return internal.Summoner{}, internal.WrapErrorf(ctx.Error, internal.ErrorCodeUnknown, "find recent")
	}

	return internal.Summoner{
		PuuId: puuid,
		AccountId: summoner.AccountId,
		SummonerId: summoner.SummonerId,
		Level: summoner.Level,
		ProfileIconId: summoner.ProfileIconId,
		Name: summoner.Name,
	}, nil
}
