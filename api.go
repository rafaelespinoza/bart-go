package bart

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const apiKey = "MW9S-E7SL-26DU-VV8V"

// baseURL is the common part of the URL for any API request. It should only be
// overridden for testing.
var baseURL = "https://api.bart.gov/api"

var httpClient = &http.Client{}

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

func requestAPI(route, cmd string, params map[string]string, out interface{}) (err error) {
	var res *http.Response

	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()

	qs := reqParams{cmd, params}
	uri := baseURL + route + "?" + qs.encode()

	res, err = httpClient.Get(uri)
	if err != nil {
		return err
	}

	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(raw, out)
	if _, ok := (err).(*json.SyntaxError); ok {
		// Handle any errors from the BART API that are formatted as XML. As of
		// this writing, the JSON API format is in beta and they never got
		// around to converting errors to JSON.
		var xmlMessage struct {
			Text    string `xml:"message>error>text"`
			Details string `xml:"message>error>details"`
		}
		xmlParseErr := xml.Unmarshal(raw, &xmlMessage)
		if xmlParseErr != nil {
			err = xmlParseErr
			return
		}
		err = fmt.Errorf("error: %s. %s", xmlMessage.Text, xmlMessage.Details)
	}

	return
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
