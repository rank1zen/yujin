package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Account struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

// Get account by riot id
// https://developer.riotgames.com/apis#account-v1/GET_getByRiotId
func (c *Client) GetAccountByRiotId(ctx context.Context, gameName, tagLine string) (*Account, error) {
	u := fmt.Sprintf(defaultAmerBaseURL+"/riot/account/v1/accounts/by-riot-id/%s/%s", gameName, tagLine)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("X-Riot-Token", os.Getenv("RIOT_API_KEY"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var m Account
	err = json.NewDecoder(res.Body).Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
