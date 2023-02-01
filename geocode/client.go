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
	ClientIDHeaderKey     = "X-NCP-APIGW-API-KEY-ID"
	ClientSecretHeaderKey = "X-NCP-APIGW-API-KEY"

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
	Kor = Lang("kor") // default
	Eng = Lang("eng") // optional
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
	clientID     *string
	clientSecret *string
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
		clientID:     &clientID,
		clientSecret: &clientSecret,
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
	if c.data != nil {
		c.data.Set(coordinate, fmt.Sprintf("%f,%f", lon, lat))
		return c
	}
	b := url.Values{}
	b.Set(coordinate, fmt.Sprintf("%f,%f", lon, lat))
	return c.copy(b)
}

// Language set language query option.
// default : "kor"
func (c *Client) Language(v Lang) *Client {
	if c.data != nil {
		c.data.Set(language, string(v))
		return c
	}
	b := url.Values{}
	b.Set(language, string(v))
	return c.copy(b)
}

// HCode set filter condition.
// you can filter array values through repetitive call.
// but if previous filter and the current filter are different,
// the previous filter will be removed.
//
// client.HCode("value1").HCode("value2)
func (c *Client) HCode(v string) *Client {
	if c.data != nil {
		f := c.data.Get(filter)
		// 필터를 설정한 적이 없으면 설정하고 리턴
		// 또는 필터를 설정한 적이 있으나 설정하려는 필터와 기존 필터가 다른 경우 기존 필터 제거하고 현재 필터 설정
		if f == "" || !strings.HasPrefix(f, string(hCode)) {
			c.data.Set(filter, fmt.Sprintf("%s@%s", hCode, v))
		} else {
			// 기존 필터에 이어서 설정
			c.data.Set(filter, fmt.Sprintf("%s;%s", f, v))
		}
		return c
	}
	b := url.Values{}
	c.data.Set(filter, fmt.Sprintf("%s@%s", hCode, v))
	return c.copy(b)
}

// BCode set filter condition.
// you can filter array values through repetitive call.
// but if previous filter and the current filter are different,
// the previous filter will be removed.
//
// client.BCode("value1").BCode("value2)
func (c *Client) BCode(v string) *Client {
	if c.data != nil {
		f := c.data.Get(filter)
		// 필터를 설정한 적이 없으면 설정하고 리턴
		// 또는 필터를 설정한 적이 있으나 설정하려는 필터와 기존 필터가 다른 경우 기존 필터 제거하고 현재 필터 설정
		if f == "" || !strings.HasPrefix(f, string(bCode)) {
			c.data.Set(filter, fmt.Sprintf("%s@%s", bCode, v))
		} else {
			// 기존 필터에 이어서 설정
			c.data.Set(filter, fmt.Sprintf("%s;%s", f, v))
		}
		return c
	}
	b := url.Values{}
	c.data.Set(filter, fmt.Sprintf("%s@%s", bCode, v))
	return c.copy(b)
}

// Page decide page you want.
// default	: 1
func (c *Client) Page(v int) *Client {
	if c.data != nil {
		c.data.Set(page, strconv.Itoa(v))
		return c
	}
	b := url.Values{}
	b.Set(page, strconv.Itoa(v))
	return c.copy(b)
}

// Count decide page unit.
// default 	: 10
// range 	: 1 ~ 100
func (c *Client) Count(v int) *Client {
	if c.data != nil {
		c.data.Set(count, strconv.Itoa(v))
		return c
	}
	b := url.Values{}
	b.Set(count, strconv.Itoa(v))
	return c.copy(b)
}

func (c *Client) copy(b url.Values) *Client {
	return &Client{
		HTTPClient:   c.HTTPClient,
		clientID:     c.clientID,
		clientSecret: c.clientSecret,
		BaseURL:      c.BaseURL,
		data:         b,
	}
}

func (c *Client) request(ctx context.Context) (*Response, error) {
	endpoint := c.BaseURL.JoinPath(Endpoint).String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = c.data.Encode()

	req.Header.Add(ClientIDHeaderKey, *c.clientID)
	req.Header.Add(ClientSecretHeaderKey, *c.clientSecret)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	var res *Response
	if err = json.Unmarshal(bytes, &res); err != nil {
		return nil, err
	}
	if res.Status != OK {
		return nil, errors.New(res.ErrorMessage)
	}
	return res, nil
}
