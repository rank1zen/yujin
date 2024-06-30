package riot

import (
	"context"
	"fmt"
)

type Summoner struct {
	AccountId     string `json:"accountId"`
	Id            string `json:"id"`
	Puuid         string `json:"puuid"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int64  `json:"summonerLevel"`
}

// Get a summoner by PUUID
//
// https://developer.riotgames.com/apis#summoner-v4/GET_getByPUUID
func (c *Client) GetSummoner(ctx context.Context, puuid string) (*Summoner, error) {
	u := fmt.Sprintf(defaultNaBaseURL+"/lol/summoner/v4/summoners/by-puuid/%v", puuid)

	req := NewRequest(WithToken2(), WithURL(u))

	summoner := new(Summoner)
	err := c.Do(ctx, req, &summoner)
	if err != nil {
		return nil, err
	}

	return summoner, nil
}
