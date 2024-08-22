package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	Progress string `json:"progess"`
	Losses   int    `json:"losses"`
	Target   int    `json:"target"`
	Wins     int    `json:"wins"`
}

// Get league entries in all queues for a given summoner ID.
// https://developer.riotgames.com/apis#league-v4/GET_getLeagueEntriesForSummoner
func (c *Client) GetLeagueEntriesForSummoner(ctx context.Context, summonerID string) (LeagueEntryList, error) {
	u := fmt.Sprintf(defaultNaBaseURL+"/lol/league/v4/entries/by-summoner/%v", summonerID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("X-Riot-Token", os.Getenv("RIOT_API_KEY"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var m LeagueEntryList
	err = json.NewDecoder(res.Body).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("decoding: %w", err)
	}

	return m, nil
}
