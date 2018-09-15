package bart

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	baseUrl = "http://api.bart.gov/api"
	apiKey  = "MW9S-E7SL-26DU-VV8V"
)

type ResponseMetaData struct {
	Uri     cDataSection
	Date    string      `json:",omitempty"`
	Time    string      `json:",omitempty"`
	Message interface{} `json:",omitempty"`
}

type cDataSection struct {
	Value string `json:"#cdata-section"`
}

func RequestApi(route, cmd string, params map[string]string, res interface{}) error {
	uri := prepareRequestUri(route, cmd, params)
	fmt.Printf("fetching %s\n", uri)

	raw, err := newGetReq(&http.Client{}, uri)
	if err != nil {
		return err
	}

	return json.Unmarshal(raw, res)
}

func prepareRequestUri(route, cmd string, params map[string]string) string {
	qs := reqParams{cmd, params}
	return baseUrl + route + "?" + qs.encode()
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
