<div class="badges">
    <a href="https://travis-ci.org/lorehov/graphite-api-client">
        <img src="https://travis-ci.org/lorehov/graphite-api-client.svg?branch=master">
    </a>
</div>


## About

Client for [graphite web api](http://graphite-api.readthedocs.io/en/latest/api.html). Currently only render API was implemented,
but metrics API will be implemented as well soon. 
Full documentation with other examples could be found on [godoc](https://godoc.org/github.com/lorehov/graphite-api-client).


## Example
```go
	package main
	
	import (
		"github.com/lorehov/graphite-api-client"
		"log"
		"fmt"
	)

	func main() {
		client, err := NewFromString("https://my-graphite.tld")
		if err != nil {
			log.Fatalf("Couldn't parse url")
		}

		res, err := client.QueryRender(
			RenderRequest{
				From: time.Unix(1468339853, 0),
				Until: time.Unix(1568339853, 0),
				MaxDataPoints: 1000,
				Targets: []string{"main"},
			},
		)

		if err != nil {
			log.Fatalf("Error %s while performing query %s", err.Error(), err.Query)
		}
		
		for _, series := range res {
			fmt.Printf("Processing target %s", series.Target)
			for _, dp := range series.Datapoints {
				fmt.Printf("Datapoint %+v", dp)
			}
		}
	}
```
