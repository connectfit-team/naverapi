package geocode

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

const (
	clientIDHeaderKey     = "X-NCP-APIGW-API-KEY-ID"
	clientSecretHeaderKey = "X-NCP-APIGW-API-KEY"

	OpenAPIBaseURL = "https://naveropenapi.apigw.ntruss.com"
	Endpoint       = "/map-geocode/v2/geocode"
)

// param key
const (
	queryConst      = "query"
	coordinateConst = "coordinate"
	filterConst     = "filter"
	languageConst   = "language"
	pageConst       = "page"
	countConst      = "count"
)

const (
	// LanguageKor korean, default
	LanguageKor = Lang("kor")
	// LanguageEng english, optional
	LanguageEng = Lang("eng") // optional
)

const (
	hCode = FilterType("HCODE")
	bCode = FilterType("BCODE")
)

// openapi response status
const (
	OK             = "OK"
	InvalidRequest = "INVALID_REQUEST"
	SystemError    = "SYSTEM_ERROR"
)

var ErrInvalidQuery = errors.New("invalid query parameter")

type Client struct {
	HTTPClient   *http.Client
	BaseURL      *url.URL
	clientID     string
	clientSecret string
}

func NewClient(
	clientID string,
	clientSecret string,
	httpClient *http.Client,
) (*Client, error) {
	baseURL, err := url.Parse(OpenAPIBaseURL)
	if err != nil {
		return nil, err
	}
	srv := &Client{
		HTTPClient:   http.DefaultClient,
		clientID:     clientID,
		clientSecret: clientSecret,
		BaseURL:      baseURL,
	}

	if httpClient != nil {
		srv.HTTPClient = httpClient
	}

	return srv, nil
}

// Query address is required value.
// Should always be called last.
func (c *Client) Query(ctx context.Context, v string, opts ...QueryOption) (*Response, error) {
	if v == "" {
		return nil, ErrInvalidQuery
	}
	data := url.Values{}
	data.Set(queryConst, v)
	for _, opt := range opts {
		opt(data)
	}
	return c.request(ctx, data)
}

func (c *Client) request(ctx context.Context, data url.Values) (*Response, error) {
	endpoint := c.BaseURL.JoinPath(Endpoint).String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = data.Encode()

	req.Header.Add(clientIDHeaderKey, c.clientID)
	req.Header.Add(clientSecretHeaderKey, c.clientSecret)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res *Response
	if err = json.Unmarshal(bytes, &res); err != nil {
		return nil, err
	}
	if res.Status != OK {
		return nil, errors.New(res.ErrorMessage)
	}
	return res, nil
}
