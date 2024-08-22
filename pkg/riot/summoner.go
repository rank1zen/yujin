package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Summoner struct {
	AccountId     string `json:"accountId"`
	Id            string `json:"id"`
	Puuid         string `json:"puuid"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int64  `json:"summonerLevel"`
}

// Get a summoner by PUUID.
//
// https://developer.riotgames.com/apis#summoner-v4/GET_getByPUUID
func (c *Client) GetSummoner(ctx context.Context, puuid string) (*Summoner, error) {
	u := fmt.Sprintf(defaultNaBaseURL+"/lol/summoner/v4/summoners/by-puuid/%v", puuid)
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
	var m *Summoner
	err = json.NewDecoder(res.Body).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("decoding: %w", err)
	}

	return m, nil
}
