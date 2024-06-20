package riot

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

const (
	apiURLFormat      = "%s://%s.%s%s"
	baseURL           = "api.riotgames.com"
	scheme            = "https"
	apiTokenHeaderKey = "X-Riot-Token"
)

func NewRiotRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Request, error) {
	apiKey := "RGAPI-aa98b358-f286-4fb7-9f11-1d93d2cf198c"
	region := "americas"

	url := fmt.Sprintf(apiURLFormat, scheme, region, baseURL, endpoint)

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add(apiTokenHeaderKey, apiKey)
	req.Header.Add("Accept", "application/json")

	return req, nil
}

type Doer interface {
	Do(request *http.Request) (*http.Response, error)
}

type Client struct {
	doer   Doer
	apiKey string
}

func NewClient(doer Doer, apiKey string) *Client {
	return &Client{
		doer:   doer,
		apiKey: apiKey,
	}
}

func (c *Client) GetMatchlist(ctx context.Context, puuid string, start int, count int) ([]string, error) {
	return listByPuuid(ctx, nil, puuid, start, count)
}

func (c *Client) GetMatch(ctx context.Context, matchId string) (*MatchDto, error) {
	return matchById(ctx, nil, matchId)
}

func (c *Client) GetMatchTimeline(ctx context.Context, matchID string) () { }

// TODO: Gets entire match list
func (c *Client) GetMatchlistFull(ctx context.Context, puuid string) (<-chan string, error) {
	ch := make(chan string)
	start, count := 0, 100

	ids, err := listByPuuid(ctx, c.doer, puuid, start, count)
	return ch, nil
}

type LeagueItem struct {
	QueueType    string      `json:"queueType"`
	SummonerName string      `json:"summonerName"`
	HotStreak    bool        `json:"hotStreak"`
	Wins         int         `json:"wins"`
	Veteran      bool        `json:"veteran"`
	Losses       int         `json:"losses"`
	FreshBlood   bool        `json:"freshBlood"`
	Inactive     bool        `json:"inactive"`
	Tier         string      `json:"tier"`
	Rank         string      `json:"rank"`
	SummonerID   string      `json:"summonerId"`
	LeaguePoints int         `json:"leaguePoints"`
}

func (c *Client) GetSummoner(ctx context.Context, puuid string) {}


func (c *Client) GetSoloqRank(ctx context.Context, puuid string) {}
