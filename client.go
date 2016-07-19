package graphite

import (
	"net/http"
	"io/ioutil"
	"time"
	"net/url"
	"errors"
)


var WrongUrlError = errors.New("Wrong url")


type RequestError struct {
	Query string
	Type string
}


func (e RequestError) Error() string {
	return e.Type
}


type Client struct {
	url url.URL
	client *http.Client
}


func NewFromString(urlString string) (*Client, error) {
	url, err := url.Parse(urlString)
	if err != nil {
		return nil, WrongUrlError
	}
	return &Client{*url, &http.Client{}}, nil
}


func (g *Client) errorResponse(r RenderRequest, t string) ([]Series, error) {
	return []Series{}, RequestError{
		Type: t,
		Query: g.queryAsString(r),
	}
}


func (g *Client) QueryRender(r RenderRequest) ([]Series, error) {
	response, err := g.client.Get(g.queryAsString(r))
	if err != nil {
		return g.errorResponse(r, "Request error")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return g.errorResponse(r, "Wrong status code")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return g.errorResponse(r, "Can't read response body")
	}

	metrics, err := unmarshallMetrics(body)
	if err != nil {
		return g.errorResponse(r, "Can't unmarshall response")
	}
	return metrics, nil
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

