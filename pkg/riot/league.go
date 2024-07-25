package riot

import (
	"context"
	"fmt"
)

type LeagueEntry struct {
	MiniSeries   *MiniSeries
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	SummonerId   string `json:"summonerId"`
	LeagueId     string `json:"leagueId"`
	Losses       int    `json:"losses"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	HotStreak    bool   `json:"hotStreak"`
	Veteran      bool   `json:"veteran"`
	FreshBlood   bool   `json:"freshBlood"`
	Inactive     bool   `json:"inactive"`
}

type LeagueEntryList []*LeagueEntry

type MiniSeries struct {
	Progress string
	Losses   int
	Target   int
	Wins     int
}

// Get league entries in all queues for a given summoner ID.
// https://developer.riotgames.com/apis#league-v4/GET_getLeagueEntriesForSummoner
func (c *Client) GetLeagueEntriesForSummoner(ctx context.Context, summonerID string) (LeagueEntryList, error) {
	u := fmt.Sprintf(defaultNaBaseURL+"/lol/league/v4/entries/by-summoner/%v", summonerID)

	req := NewRequest(WithToken2(), WithURL(u))

	var a LeagueEntryList
	err := c.Do(ctx, req, &a)
	if err != nil {
		return nil, err
	}

	return a, nil
}
