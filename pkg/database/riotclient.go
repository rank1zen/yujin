package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/rank1zen/yujin/pkg/logging"
)

type RiotClient interface {
	GetSummoner(puuid string) (*lol.Summoner, error)
	GetMatchlist(puuid string) ([]string, error)
	GetMatch(matchId string) (*lol.Match, error)
	GetLeagueBySummoner(summonerId string)  (*lol.LeagueItem, error)
}

type golioClient struct {
	golio *golio.Client
}

func NewGolioClient(ctx context.Context, apiKey string) RiotClient {
	log := logging.FromContext(ctx).Sugar()

	log.Infof("starting golio client...")
	return &golioClient{
		golio: golio.NewClient(apiKey, golio.WithRegion(api.RegionNorthAmerica)),
	}
}

func WithNewGolioClient(ctx context.Context, e interface{ SetRiotClient(RiotClient) }, apiKey string) {
	gc := NewGolioClient(ctx, apiKey)
	e.SetRiotClient(gc)
}

func (g *golioClient) GetSummoner(puuid string) (*lol.Summoner, error) {
	return g.golio.Riot.LoL.Summoner.GetByPUUID(puuid)
}

func (g *golioClient) GetMatch(matchId string) (*lol.Match, error) {
	return g.golio.Riot.LoL.Match.Get(matchId)
}

func (g *golioClient) GetMatchlist(matchId string) ([]string, error) {
	return nil, fmt.Errorf("no implemented")
}

func (g *golioClient) GetLeagueBySummoner(summonerId string) (*lol.LeagueItem, error) {
	leagues, err := g.golio.Riot.LoL.League.ListBySummoner(summonerId)
	if err != nil {
		return nil, fmt.Errorf("riot api: %w", err)
	}

	for _, l := range leagues {
		if l.QueueType == "FIXME" {
			return l, nil
		}
	}
	return nil, errors.New("soloq not found")
}
