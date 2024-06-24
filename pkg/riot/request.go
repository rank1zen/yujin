package riot

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const userAgent = "server:yujin:v1.0"

type Request struct {
	body               url.Values
	query              url.Values
	method             string
	token              string
	url                string
	auth               string
	tags               []string
	emptyResponseBytes int
	retry              bool
	client             *http.Client
}

type RequestOption func(*Request)

func NewRequest(opts ...RequestOption) *Request {
	req := &Request{
		body:   url.Values{},
		query:  url.Values{},
		method: "GET",
		url:    "",

		token: "",
		auth:  "",

		tags: nil,

		emptyResponseBytes: 0,
		retry:              true,
		client:             nil,
	}

	for _, opt := range opts {
		opt(req)
	}

	return req
}

func (r *Request) HTTPRequest(ctx context.Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, r.method, r.url, strings.NewReader(r.body.Encode()))
	if err != nil {
		return nil, err
	}

	// TODO: this is kinda bs bruv
	// req.URL.RawQuery = r.query.Encode()

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", userAgent)

	if r.token != "" {
		req.Header.Add("X-Riot-Token", r.token)
	}

	return req, err
}

func WithToken2() RequestOption {
	return WithToken(os.Getenv("RIOT_API_KEY"))
}

func WithToken(token string) RequestOption {
	return func(req *Request) {
		req.token = token
	}
}

func WithQuery(key, val string) RequestOption {
	if val == "" {
		return func(req *Request) {}
	}

	return func(req *Request) {
		req.query.Set(key, val)
	}
}

func WithURL(url string) RequestOption {
	return func(req *Request) {
		req.url = url
	}
}
