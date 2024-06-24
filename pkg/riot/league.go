package riot

import "context"

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
func (c *Client) GetSoloqRank(ctx context.Context, summonerID string) (*LeagueItem, error) {
	return nil, nil
}
