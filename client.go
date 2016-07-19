package graphite

import (
	"net/http"
	"io/ioutil"
	"time"
	"net/url"
	"errors"
)

type Client struct {
	url url.URL
	client *http.Client
}


func NewFromString(urlString string) (*Client, error) {
	url, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	return &Client{*url, &http.Client{}}, nil
}


// TODO: here we must handle different errors and expose our errors
func (g *Client) QueryRender(r RenderRequest) ([]Series, error) {
	response, err := g.client.Get(g.queryAsString(r))
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Wrong response code")
	}

	body, err := ioutil.ReadAll(response.Body)
	metrics, err := unmarshallMetrics(body)
	return metrics, err
}


func (g *Client) QueryRenderFromUntil(from, until time.Time, targets []string) ([]Series, error) {
	return g.QueryRender(
		RenderRequest{
			From: from,
			Until: until,
			Targets: targets,
		})
}


func (g *Client) queryAsString(r RenderRequest) string {
	return g.url.String() + r.ToQueryString()
}

