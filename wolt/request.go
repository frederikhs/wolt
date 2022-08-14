package wolt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL    *url.URL
	HttpClient *http.Client
}

type RoundTripper struct {
	r     http.RoundTripper
	token string
}

func NewClient(token string) *Client {
	return &Client{
		BaseURL: &url.URL{
			Scheme: "https",
			Host:   "restaurant-api.wolt.com",
		},
		HttpClient: &http.Client{
			Timeout:   time.Second * 10,
			Transport: RoundTripper{r: http.DefaultTransport, token: token},
		},
	}
}

func (c *Client) RequestOrders(limit, skip int) (*[]FullOrder, error) {
	endpointUrl := c.constructUrl("/v2/order_details/")

	q := url.Values{}
	q.Add("limit", fmt.Sprintf("%d", limit))
	q.Add("skip", fmt.Sprintf("%d", skip))
	endpointUrl.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, endpointUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	var orders []FullOrder
	return &orders, c.handleRequestResponse(req, &orders)
}

func (c *Client) constructUrl(path string) *url.URL {
	return c.BaseURL.ResolveReference(&url.URL{Path: path})
}

func (mrt RoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("Authorization", "Bearer "+mrt.token)
	return mrt.r.RoundTrip(r)
}

func (c *Client) handleRequestResponse(r *http.Request, i interface{}) error {
	res, err := c.HttpClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("got http status code: %d", res.StatusCode)
	}

	err = json.Unmarshal(body, i)
	if err != nil {
		return err
	}

	return nil
}
