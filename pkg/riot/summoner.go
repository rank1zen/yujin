package riot

import (
	"context"
	"fmt"
)

type Summoner struct {
	AccountId     *string `json:"accountId"`
	ProfileIconId *int    `json:"profileIconId"`
	RevisionDate  *int64  `json:"revisionDate"`
	Id            *string `json:"id"`
	Puuid         *string `json:"puuid"`
	SummonerLevel *int64  `json:"summonerLevel"`
}

// Get a summoner by PUUID.
// TODO: test this please
func (c *Client) GetSummoner(ctx context.Context, puuid string) (*Summoner, error) {
	_ = fmt.Sprintf("/lol/summoner/v4/summoners/by-puuid/%v", puuid)
	return nil, nil
}
