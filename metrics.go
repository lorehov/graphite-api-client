package graphite

import (
	"time"
	"net/url"
	"strconv"
)

type FindMetricRequest struct {
	From	time.Time
	Until	time.Time
	Query string
	Wildcards bool
}


func (r FindMetricRequest) toQueryString() string {
	values := url.Values{
		"format": []string{"treejson"},
	}
	if !r.From.IsZero() {
		values.Set("from", strconv.FormatInt(r.From.Unix(), 10))
	}
	if !r.Until.IsZero() {
		values.Set("until", strconv.FormatInt(r.Until.Unix(), 10))
	}
	if r.Query != "" {
		values.Set("query", r.Query)
	}
	if r.Wildcards {
		values.Set("wildcards", strconv.Itoa(1))
	}
	qs := values.Encode()
	return "/metrics/find?" + qs
}

//{"text": "geo", "expandable": 1, "leaf": 0, "id": "cluster.geo", "allowChildren": 1}
type Metric struct {
	Id string
	Text string
	Expandable int
	Leaf int
	AllowChildren int
}


func (g *Client) FindMetrics(r FindMetricRequest) ([]Metric, error) {
	return nil, nil
}


type ExpandMetricRequest struct {
	Query string
	GroupByExpr bool
	LeavesOnly bool
}


func (r ExpandMetricRequest) toQueryString() string {
	values := url.Values{}
	if r.Query != "" {
		values.Set("query", r.Query)
	}
	if r.GroupByExpr {
		values.Set("groupByExpr", strconv.Itoa(1))
	}
	if r.LeavesOnly {
		values.Set("leavesOnly", strconv.Itoa(1))
	}
	qs := values.Encode()
	return "/metrics/expand?" + qs
}


type ExpandResult struct {
	Results []string
}


//{"results": ["cluster.geo", "cluster.mail"]}
func (g *Client) ExpandMetrics(r ExpandMetricRequest) (ExpandResult, error) {
	return ExpandResult{}, nil
}
