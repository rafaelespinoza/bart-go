package bart

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	baseURL = "https://api.bart.gov/api"
	apiKey  = "MW9S-E7SL-26DU-VV8V"
)

// Client gives you easy access to several BART API endpoints. See examples for
// general usage.
type Client struct {
	*AdvisoriesAPI
	*EstimatesAPI
	*RoutesAPI
	*SchedulesAPI
	*StationsAPI
}

// ResponseMetaData is contains some data about the response. Not all of the
// fields are filled by every API endpoint.
type ResponseMetaData struct {
	URI     CDATASection
	Date    string      `json:",omitempty"`
	Time    string      `json:",omitempty"`
	Message interface{} `json:",omitempty"`
}

// CDATASection is merely a helper for unmarshaling certain fields. The original
// BART API has long returned XML instead of JSON, and its presence is an
// artifact of BART's output conversion. This type is meant for internal use and
// is only exported for documentation purposes since it shows up in so many
// other type definitions in this package.
type CDATASection struct {
	Value string `json:"#cdata-section"`
}

func requestAPI(route, cmd string, params map[string]string, res interface{}) error {
	uri := prepareRequestURI(route, cmd, params)
	fmt.Printf("fetching %s\n", uri)

	raw, err := newGetReq(&http.Client{}, uri)
	if err != nil {
		return err
	}

	return json.Unmarshal(raw, res)
}

func prepareRequestURI(route, cmd string, params map[string]string) string {
	qs := reqParams{cmd, params}
	return baseURL + route + "?" + qs.encode()
}

func newGetReq(c *http.Client, uri string) ([]byte, error) {
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return []byte{}, err
	}

	res, err := c.Do(req)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	return body, err
}

type reqParams struct {
	cmd   string
	pairs map[string]string
}

func (p reqParams) encode() string {
	qs := url.Values{}
	qs.Add("cmd", p.cmd)
	qs.Add("json", "y")
	qs.Add("key", apiKey)

	for k, v := range p.pairs {
		qs.Add(k, v)
	}

	return qs.Encode()
}
