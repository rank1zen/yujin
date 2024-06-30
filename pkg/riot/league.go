package riot

import (
	"context"
	"fmt"
)

type LeagueItem struct {
	QueueType    string `json:"queueType"`
	SummonerName string `json:"summonerName"`
	HotStreak    bool   `json:"hotStreak"`
	Wins         int    `json:"wins"`
	Veteran      bool   `json:"veteran"`
	Losses       int    `json:"losses"`
	FreshBlood   bool   `json:"freshBlood"`
	Inactive     bool   `json:"inactive"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	SummonerID   string `json:"summonerId"`
	LeaguePoints int    `json:"leaguePoints"`
}

// Get ranked soloq for a given summoner ID.
// TODO: implement
// https://developer.riotgames.com/apis#league-v4/GET_getLeagueEntriesForSummoner
func (c *Client) GetSoloqRank(ctx context.Context, summonerID string) (*LeagueItem, error) {
	panic("not implemented")
	u := fmt.Sprintf(defaultNaBaseURL+"/lol/league/v4/entries/by-summoner/%v", summonerID)

	req := NewRequest(WithToken2(), WithURL(u))

	var a []LeagueItem
	err := c.Do(ctx, req, &a)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
