// Package graphite provides client to graphite-web-api (http://graphite-api.readthedocs.io/en/latest/api.html).
package graphite

import (
	"net/http"
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


func (c *Client) errorResponse(r qsGenerator, t string) ([]Series, error) {
	return []Series{}, RequestError{
		Type: t,
		Query: c.queryAsString(r),
	}
}


func (c *Client) queryAsString(r qsGenerator) string {
	return c.Url.String() + r.toQueryString()
}


type qsGenerator interface {
	toQueryString() string
}

