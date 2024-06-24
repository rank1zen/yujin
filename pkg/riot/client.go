package riot

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rank1zen/yujin/pkg/logging"
)

const (
	defaultNaBaseURL       = "https://na1.api.riotgames.com"
	defaultAmerBaseURL     = "https://americas.api.riotgames.com"
	defaultBaseURLTemplate = "https://%s.api.riotgames.com"
)

type Client struct {
	client      *http.Client
	defaultOpts []RequestOption // TODO: we havent implemented this feat yet 
}

func NewClient(opts ...RequestOption) *Client {
	httpClient := &http.Client{
		Timeout: 4 * time.Second,
	}

	return &Client{
		client:      httpClient,
		defaultOpts: opts,
	}
}

func (c *Client) doRequest(ctx context.Context, req *Request) (*http.Response, error) {
	logger := logging.FromContext(ctx).Sugar()

	rq, err := req.HTTPRequest(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(rq)
	if err != nil {
		select {
		case <-ctx.Done():
			logger.Debugf("context cancelled: %v", ctx.Err())
			return nil, ctx.Err()
		default:
		}

		logger.Debugf("interal client: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logger.Debugf("http: %v at %v", resp.Status, req.url)

		err, found := StatusToError[resp.StatusCode]
		if !found {
			return nil, Error{Message: "unknown status", StatusCode: resp.StatusCode}
		}

		return nil, err
	}

	return resp, nil
}

func (c *Client) Do(ctx context.Context, req *Request, dst any) error {
	logger := logging.FromContext(ctx).Sugar()

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(dst)
	if err != nil {
		logger.Debugf("failed to decode: %v", err)
		return err
	}

	return err
}
