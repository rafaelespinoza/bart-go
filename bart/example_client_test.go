package bart_test

import (
	"fmt"
	"net/http"

	"github.com/rafaelespinoza/bart-go/bart"
)

type demo struct{}

func (*demo) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Println(req.URL)
	return http.DefaultTransport.RoundTrip(req)
}

func ExampleClient() {
	var client *bart.Client

	// use default settings
	client = bart.NewClient(nil)

	// or configure client settings
	conf := &bart.Config{
		// replace with your registered API key, optional.
		Key: "FOO-BAR",

		// configure HTTP client, optional.
		HTTP: &http.Client{
			Transport: &demo{},
		},
	}
	client = bart.NewClient(conf)

	res, err := client.RequestStations()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}
