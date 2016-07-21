package graphite

import (
	"net/http"
	"io/ioutil"
	"net/url"
	"errors"
)


var WrongUrlError = errors.New("Wrong url")


// RequestError is special error, which not only implements an error interface,
// but also provides access to `Query` field, containing original query which
// cause an error.
type RequestError struct {
	Query string
	Type string
}


func (e RequestError) Error() string {
	return e.Type
}


// Client is client to `graphite web api` (http://graphite-api.readthedocs.io/en/latest/api.html).
// You can either instantiate it manually by providing `url` and `client` or use a `NewFromString` shortcut.
type Client struct {
	Url url.URL
	Client *http.Client
}


// NewFromString is a convenient function for constructing `graphite.Client` from url string.
// For example 'https://my-graphite.tld'. NewFromString could return either `graphite.Client`
// instance or `WrongUrlError` error.
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


// QueryRender performs query to graphite `/render/` api. Normally it should return `[]graphite.Series`,
// but if things go wrong it will return `graphite.RequestError` error.
func (g *Client) QueryRender(r RenderRequest) ([]Series, error) {
	response, err := g.Client.Get(g.queryAsString(r))
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


func (g *Client) queryAsString(r RenderRequest) string {
	return g.Url.String() + requestToQueryString(r)
}

