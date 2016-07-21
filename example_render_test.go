package graphite

import (
	"net/http"
	"net/http/httptest"
	"fmt"
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
			Targets: []string{"main"},
		},
	)
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
