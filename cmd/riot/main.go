package main

import (
	"os"

	"github.com/KnutZuidema/golio"
	"github.com/rank1zen/yujin/internal/riot"
)

func main() {
	client := golio.NewClient(os.Getenv("API_KEY"))
	riot.NewSummonerQ(client.Riot.LoL.Summoner)
}
