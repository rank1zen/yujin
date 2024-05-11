package database

import (
	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
)

func NewGolioClient(apiKey string) *golio.Client {
        return golio.NewClient(apiKey, golio.WithRegion(api.RegionNorthAmerica))
}

func WithNewGolioClient(apiKey string, e interface{ SetGolioClient(*golio.Client) }) {
        gc := NewGolioClient(apiKey)
        e.SetGolioClient(gc)
}
