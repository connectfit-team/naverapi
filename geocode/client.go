package geocode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	clientIDHeaderKey     = "X-NCP-APIGW-API-KEY-ID"
	clientSecretHeaderKey = "X-NCP-APIGW-API-KEY"

	OpenAPIBaseURL = "https://naveropenapi.apigw.ntruss.com"
	Endpoint       = "/map-geocode/v2/geocode"
)

// param key
const (
	query      = "query"
	coordinate = "coordinate"
	filter     = "filter"
	language   = "language"
	page       = "page"
	count      = "count"
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
	data         url.Values
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
func (c *Client) Query(ctx context.Context, v string) (*Response, error) {
	if v == "" {
		return nil, ErrInvalidQuery
	}
	if c.data == nil {
		c.data = url.Values{}
	}
	c.data.Set(query, v)
	return c.request(ctx)
}

// Coordinate set coordinates to be the center of the search.
// If set, computes the distance from the `Query()` value to the coordinates.
func (c *Client) Coordinate(lon, lat float64) *Client {
	cc := c.clone()
	if cc.data == nil {
		cc.data = url.Values{}
	}
	cc.data.Set(coordinate, fmt.Sprintf("%f,%f", lon, lat))
	return cc
}

// Language set language query option.
// default : "kor"
func (c *Client) Language(v Lang) *Client {
	cc := c.clone()
	if cc.data == nil {
		cc.data = url.Values{}
	}
	cc.data.Set(language, string(v))
	return cc
}

// HCode set filter condition.
// you can filter array values through repetitive call.
// but if previous filter and the current filter are different,
// the previous filter will be removed.
func (c *Client) HCode(v ...string) *Client {
	cc := c.clone()
	vj := strings.Join(v, ";")
	if cc.data == nil {
		cc.data = url.Values{}
		cc.data.Set(filter, fmt.Sprintf("%s@%s", hCode, vj))
	} else {
		f := cc.data.Get(filter)
		// 필터를 설정한 적이 없으면 설정하고 리턴
		// 또는 필터를 설정한 적이 있으나 설정하려는 필터와 기존 필터가 다른 경우 기존 필터 제거하고 현재 필터 설정
		if f == "" || !strings.HasPrefix(f, string(hCode)) {
			cc.data.Set(filter, fmt.Sprintf("%s@%s", hCode, vj))
		} else {
			// 기존 필터에 이어서 설정
			cc.data.Set(filter, fmt.Sprintf("%s;%s", f, vj))
		}
	}
	return cc
}

// BCode set filter condition.
// you can filter array values through repetitive call.
// but if previous filter and the current filter are different,
// the previous filter will be removed.
func (c *Client) BCode(v ...string) *Client {
	cc := c.clone()
	vj := strings.Join(v, ";")
	if cc.data == nil {
		cc.data = url.Values{}
		cc.data.Set(filter, fmt.Sprintf("%s@%s", bCode, vj))
	} else {
		f := cc.data.Get(filter)
		// 필터를 설정한 적이 없으면 설정하고 리턴
		// 또는 필터를 설정한 적이 있으나 설정하려는 필터와 기존 필터가 다른 경우 기존 필터 제거하고 현재 필터 설정
		if f == "" || !strings.HasPrefix(f, string(bCode)) {
			cc.data.Set(filter, fmt.Sprintf("%s@%s", bCode, vj))
		} else {
			// 기존 필터에 이어서 설정
			cc.data.Set(filter, fmt.Sprintf("%s;%s", f, vj))
		}
	}
	return cc
}

// Page decide page you want.
// default	: 1
func (c *Client) Page(v int) *Client {
	cc := c.clone()
	if cc.data == nil {
		cc.data = url.Values{}
	}
	cc.data.Set(page, strconv.Itoa(v))
	return cc
}

// Count decide page unit.
// default 	: 10
// range 	: 1 ~ 100
func (c *Client) Count(v int) *Client {
	cc := c.clone()
	if cc.data == nil {
		cc.data = url.Values{}
	}
	cc.data.Set(count, strconv.Itoa(v))
	return cc
}

func (c *Client) clone() *Client {
	cc := *c
	return &cc
}

func (c *Client) request(ctx context.Context) (*Response, error) {
	endpoint := c.BaseURL.JoinPath(Endpoint).String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = c.data.Encode()

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
