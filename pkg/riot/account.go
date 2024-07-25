package riot

import (
	"context"
	"fmt"
)

type Account struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

// Get account by riot id
// https://developer.riotgames.com/apis#account-v1/GET_getByRiotId
func (c *Client) GetAccountByRiotId(ctx context.Context, gameName, tagLine string) (*Account, error) {
	u := fmt.Sprintf(defaultNaBaseURL+"/account/v1/accounts/by-riot-id/%v/%v", gameName, tagLine)

	req := NewRequest(WithToken2(), WithURL(u))

	var a Account
	err := c.Do(ctx, req, &a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
