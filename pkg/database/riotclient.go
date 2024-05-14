package database

import (
	"context"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/rank1zen/yujin/pkg/logging"
)

type RiotClient interface {
	GetSummoner(puuid string) (*lol.Summoner, error)
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
