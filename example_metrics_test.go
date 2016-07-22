package graphite

import (
	"net/http/httptest"
	"net/http"
	"fmt"
)

func createFindMetricsTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[
			{"text": "one", "expandable": 1, "leaf": 0, "id": "cluster.one", "allowChildren": 1},
			{"text": "two", "expandable": 1, "leaf": 0, "id": "cluster.two", "allowChildren": 1}]`)
	}))
}


func ExampleClient_FindMetrics() {
	ts := createFindMetricsTestServer()
	defer ts.Close()

	client, _ := NewFromString(ts.URL)
	res, _ := client.FindMetrics(
		FindMetricRequest{
			Query: "cluster.*",
		},
	)
	fmt.Printf("%+v", res)
	// Output: [{Id:cluster.one Text:one Expandable:1 Leaf:0 AllowChildren:1} {Id:cluster.two Text:two Expandable:1 Leaf:0 AllowChildren:1}]

}


func createExpandMetricsTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{\"results\": [\"cluster.one\", \"cluster.two\"]}")
	}))
}

func ExampleClient_ExpandMetrics() {
	ts := createExpandMetricsTestServer()
	defer ts.Close()

	client, _ := NewFromString(ts.URL)
	res, _ := client.ExpandMetrics(
		ExpandMetricRequest{
			Query: "cluster.*",
		},
	)
	fmt.Printf("%+v", res)
	// Output: {Results:[cluster.one cluster.two]}
}
