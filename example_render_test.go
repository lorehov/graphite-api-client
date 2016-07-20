package graphite

import (
	"net/http"
	"net/http/httptest"
	"fmt"
	"time"
	"strconv"
)

func createTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "[{\"target\": \"main\", \"datapoints\": [[1, 1468339853], [1.0, 1468339854], [null, 1468339855]]}]")
	}))
}


func ExampleClient_QueryRender() {
	ts := createTestServer()
	defer ts.Close()

	client, _ := NewFromString(ts.URL)
	res, _ := client.QueryRender(
		RenderRequest{
			From: time.Unix(1468339853, 0),
			Until: time.Unix(1568339853, 0),
			MaxDataPoints: 1000,
			Targets: []string{"main"},
		},
	)
	fmt.Printf(strconv.FormatInt(res[0].Datapoints[0].Timestamp.Unix(), 10))
	// Output: 1468339853
}


func ExampleClient_QueryRenderFromUntil() {
	ts := createTestServer()
	defer ts.Close()

	client, _ := NewFromString(ts.URL)
	res, _ := client.QueryRenderFromUntil(time.Unix(1468339853, 0), time.Unix(1568339853, 0), []string{"main"})
	fmt.Printf(strconv.FormatInt(res[0].Datapoints[0].Timestamp.Unix(), 10))
	// Output: 1468339853

}

func ExampleNewFromString() {
	client, _ := NewFromString("http://my-graphite.tld")
	fmt.Printf(client.Url.Host)
	// Output: my-graphite.tld
}

func ExampleRequestError_Error() {
	err := RequestError{
		Type: "Error type",
	}
	fmt.Printf(err.Error())
	// Output: Error type
}
